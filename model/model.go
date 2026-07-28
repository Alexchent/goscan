package model

import (
	"github.com/Alexchent/goscan/help"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// MyFile 记录文件详情（按需计算 MD5），与 ScanRecord 共存。
// ScanRecord 用于扫描去重，MyFile 用于重复文件分析。
type MyFile struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	FullFileName string `gorm:"uniqueIndex"`
	Size         int64  // 文件大小
	FileMD5      string
	Created      time.Time
	Status       int
}

// ScanRecord 扫描去重记录表，存原始 "path,size"。
// 用 UNIQUE(full_file_name) 约束 + ON CONFLICT DO NOTHING 实现幂等插入，
// 替换原 Redis Set 的去重语义。
type ScanRecord struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	FullFileName string `gorm:"uniqueIndex"`
	Size         int64
	Created      time.Time
}

// recordEntry 是 ScanRecord 的一条导出表示，供 logic/cmd 使用
type RecordEntry struct {
	Path string
	Size int64
}

type IScanFile interface {
	Insert(fullFileName string)
	FindByFullFileName(fullFileName string) *MyFile
	FindMd5(fileMD5 string) []MyFile
}

type ScanFile struct {
	db *gorm.DB
}

// 接口完整性校验
var _ IScanFile = (*ScanFile)(nil)

// NewScanFile 打开/创建 SQLite 数据库（scan.db）并迁移 schema。
// dsn 为目录路径，数据库文件落在该目录下。
func NewScanFile(dsn string) *ScanFile {
	if len(dsn) > 0 && dsn[len(dsn)-1] != '/' {
		dsn += "/"
	}
	db, err := gorm.Open(sqlite.Open(dsn+"scan.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 启用 WAL，提升并发写入
	db.Exec("PRAGMA journal_mode=WAL")
	// 迁移 schema
	if err := db.AutoMigrate(&MyFile{}, &ScanRecord{}); err != nil {
		panic(err)
	}
	return &ScanFile{db: db}
}

// DB 暴露底层 gorm.DB，供 logic 层组合查询使用
func (m *ScanFile) DB() *gorm.DB {
	return m.db
}

func (m *ScanFile) Insert(fullFileName string) {
	var fileMD5 string
	var fileSize int64
	result := m.FindByFullFileName(fullFileName)
	stat, err := os.Stat(fullFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("文件 %s 不存在\n", fullFileName)
		}
		fileSize = 0
	} else {
		fileSize = stat.Size()
		fileMD5 = help.GetFileMd5(fullFileName)
	}
	if result.FullFileName == "" {
		m.db.Create(&MyFile{
			FullFileName: fullFileName,
			Size:         fileSize,
			FileMD5:      fileMD5,
			Created:      time.Now(),
			Status:       1,
		})
	} else if result.Size == 0 {
		result.Size = fileSize
		m.db.Save(result)
	}
}

func (m *ScanFile) FindByFullFileName(fullFileName string) *MyFile {
	var data MyFile
	m.db.First(&data, "full_file_name = ?", fullFileName)
	return &data
}

func (m *ScanFile) FindMd5(fileMD5 string) []MyFile {
	var data []MyFile
	m.db.Find(&data, "file_md5 = ?", fileMD5)
	return data
}

// ---- ScanRecord 操作（替换原 Redis Set 链路）----

// AddRecord 插入一条扫描记录，已存在则忽略（去重）。
// 返回 true 表示是新插入的文件，false 表示已存在（等价于原 cache.SAdd==1）。
func (m *ScanFile) AddRecord(path string, size int64) bool {
	rec := ScanRecord{FullFileName: path, Size: size, Created: time.Now()}
	// sqlite: INSERT OR IGNORE
	res := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rec)
	return res.RowsAffected == 1
}

// AllRecords 返回全部扫描记录（等价于 cache.SMembers）。
// 为兼容旧日志格式每条返回 "path,size" 字符串。
func (m *ScanFile) AllRecords() []string {
	var rows []ScanRecord
	// 按 path 升序，结果更稳定
	m.db.Order("full_file_name ASC").Find(&rows)
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.FullFileName+","+strconv.FormatInt(r.Size, 10))
	}
	return out
}

// AllEntries 返回结构化的全部记录，供需要 path/size 的调用方使用
func (m *ScanFile) AllEntries() []RecordEntry {
	var rows []ScanRecord
	m.db.Order("full_file_name ASC").Find(&rows)
	out := make([]RecordEntry, 0, len(rows))
	for _, r := range rows {
		out = append(out, RecordEntry{Path: r.FullFileName, Size: r.Size})
	}
	return out
}

// Search 模糊查找（等价于原 SearchFromRedisSet）。
// 过滤掉 - 和 _ 后做大小写不敏感包含匹配，返回 "path,size" 字符串。
func (m *ScanFile) Search(keyword string) []string {
	reg := strings.NewReplacer("-", "", "_", "")
	// 还原区间: keyword 用 sql LIKE 无法表达"去掉 - 和 _ 再匹配"，
	// 因此先取全量再内存过滤，保持与原 Redis 版本一致的匹配语义。
	kw := strings.ToLower(reg.Replace(keyword))
	if kw == "" {
		return nil
	}
	var rows []ScanRecord
	m.db.Order("full_file_name ASC").Find(&rows)
	out := make([]string, 0)
	for _, r := range rows {
		if strings.Contains(strings.ToLower(reg.Replace(r.FullFileName)), kw) {
			out = append(out, r.FullFileName+","+strconv.FormatInt(r.Size, 10))
		}
	}
	return out
}

// DeleteBySuffix 删除 fullFileName 以 suffix 结尾的记录（等价于按后缀清理）
func (m *ScanFile) DeleteBySuffix(suffix string) (paths []string) {
	if suffix == "" {
		return nil
	}
	var rows []ScanRecord
	m.db.Where("full_file_name LIKE ?", "%"+suffix).Find(&rows)
	for _, r := range rows {
		m.db.Delete(&r)
		paths = append(paths, r.FullFileName+","+strconv.FormatInt(r.Size, 10))
	}
	return paths
}

// DeleteByKeyword 删除 fullFileName 包含 keyword 的记录
func (m *ScanFile) DeleteByKeyword(keyword string) (paths []string) {
	if keyword == "" {
		return nil
	}
	var rows []ScanRecord
	m.db.Where("full_file_name LIKE ?", "%"+keyword+"%").Find(&rows)
	for _, r := range rows {
		m.db.Delete(&r)
		paths = append(paths, r.FullFileName+","+strconv.FormatInt(r.Size, 10))
	}
	return paths
}

// DeleteByFromLogRecord 按字符串输入删除记录（用于 load 反向导入场景之外，
// clean 之外罕见场景）。这里提供按 "path" 精确删除的便捷方法。
func (m *ScanFile) DeleteByPath(path string) bool {
	res := m.db.Where("full_file_name = ?", path).Delete(&ScanRecord{})
	return res.RowsAffected > 0
}

// Count 返回记录总数
func (m *ScanFile) Count() int64 {
	var n int64
	m.db.Model(&ScanRecord{}).Count(&n)
	return n
}
