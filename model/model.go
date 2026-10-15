package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// ScanRecord 扫描去重记录表，存原始 "path,size"。
// 用 UNIQUE(full_file_name) 约束 + ON CONFLICT DO NOTHING 实现幂等插入，
// 替换原 Redis Set 的去重语义。
// NormFileName 是 full_file_name 归一化后的副本（小写、去掉 `-`/`_`），
// 供 Search 走 SQL `LIKE` 而非全表内存过滤。
type ScanRecord struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	FullFileName string `gorm:"uniqueIndex"`
	NormFileName string `gorm:"index"`
	Size         int64
	Created      time.Time
}

type IScanFile interface {
	AddRecord(path string, size int64) bool
}

type ScanFile struct {
	db *gorm.DB
}

// 接口完整性校验
var _ IScanFile = (*ScanFile)(nil)

// parseLogLevel 将配置文件中的字符串级别解析为 gorm logger.LogLevel。
// 取值（大小写不敏感）：silent / error / warn / info。
// 空字符串或无法识别时默认 silent，避免打印慢 SQL 等噪声日志。
func parseLogLevel(s string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	case "silent", "":
		return logger.Silent
	default:
		return logger.Silent
	}
}

// NewScanFile 打开/创建 SQLite 数据库（scan.db）并迁移 schema。
// dsn 为目录路径，数据库文件落在该目录下。
// logLevel 控制 GORM SQL 日志级别，取值见 parseLogLevel：silent/error/warn/info。
// 传空字符串时默认 silent（不输出任何 SQL 日志）。
func NewScanFile(dsn, logLevel string) *ScanFile {
	if len(dsn) > 0 && dsn[len(dsn)-1] != '/' {
		dsn += "/"
	}
	db, err := gorm.Open(sqlite.Open(dsn+"scan.db"), &gorm.Config{
		Logger: logger.Default.LogMode(parseLogLevel(logLevel)),
	})
	if err != nil {
		panic("failed to connect database")
	}
	// 启用 WAL，提升并发写入
	db.Exec("PRAGMA journal_mode=WAL")
	// 迁移 schema
	if err := db.AutoMigrate(&ScanRecord{}); err != nil {
		panic(err)
	}
	sf := &ScanFile{db: db}
	// 为存量数据回填 NormFileName（首次升级时执行）
	sf.backfillNorm()
	return sf
}

// backfillNorm 为 norm_file_name 为空的历史行回填归一化值。
// 幂等：仅处理空值行；新写入的路径在 AddRecord/AddRecords 时已填好。
func (m *ScanFile) backfillNorm() {
	var rows []ScanRecord
	m.db.Where("norm_file_name = '' OR norm_file_name IS NULL").Find(&rows)
	if len(rows) == 0 {
		return
	}
	var fullNames = make([]string, 0, len(rows))
	for _, r := range rows {
		fullNames = append(fullNames, r.FullFileName)
	}
	// 逐条 UPDATE 回填（存量回填一次性，量通常有限；如需极致可改为批量 CASE 更新）
	for _, name := range fullNames {
		m.db.Model(&ScanRecord{}).Where("full_file_name = ?", name).
			Update("norm_file_name", normalizeName(name))
	}
}

// DB 暴露底层 gorm.DB，供 logic 层组合查询使用
func (m *ScanFile) DB() *gorm.DB {
	return m.db
}

// ---- ScanRecord 操作（替换原 Redis Set 链路）----

// normalizeName 生成用于模糊匹配的归一化名：小写并去掉 `-` 和 `_`。
// 与历史 Redis Go 版 SearchFromRedisSet 的匹配语义保持一致。
func normalizeName(name string) string {
	reg := getReplacer()
	return strings.ToLower(reg.Replace(name))
}

// getReplacer 返回共用的 -/_ 替换器（strings.Replacer 非并发安全，多 goroutine 需加锁，
// 但这里调用路径均为单 goroutine 扫描/查询，直接复用安全）。
func getReplacer() *strings.Replacer {
	return strings.NewReplacer("-", "", "_", "")
}

// AddRecord 插入一条扫描记录，已存在则忽略（去重）。
// 返回 true 表示是新插入的文件，false 表示已存在（等价于原 cache.SAdd==1）。
func (m *ScanFile) AddRecord(path string, size int64) bool {
	rec := ScanRecord{
		FullFileName: path,
		NormFileName: normalizeName(path),
		Size:         size,
		Created:      time.Now(),
	}
	// sqlite: INSERT OR IGNORE
	res := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rec)
	return res.RowsAffected == 1
}

// AddRecords 批量插入多条扫描记录，用单个事务 + INSERT OR IGNORE 一次提交，
// 避免扫描大量文件时逐条磁盘 IO。返回本批中实际新增（此前不存在）的记录切片，
// 供调用方精确打印与写日志。
func (m *ScanFile) AddRecords(recs []*ScanRecord) (added []ScanRecord) {
	if len(recs) == 0 {
		return nil
	}
	fullNames := make([]string, 0, len(recs))
	for _, r := range recs {
		fullNames = append(fullNames, r.FullFileName)
	}

	// 批前已存在路径集合（用于筛选新增记录，单进程场景下无并发竞争）
	preExist, err := m.existing(fullNames)
	if err != nil {
		return nil
	}

	now := time.Now()
	batch := make([]ScanRecord, 0, len(recs))
	for _, r := range recs {
		batch = append(batch, ScanRecord{
			FullFileName: r.FullFileName,
			NormFileName: normalizeName(r.FullFileName),
			Size:         r.Size,
			Created:      now,
		})
	}
	if err := m.db.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error
	}); err != nil {
		return nil
	}

	// 筛选其中新增（不在 preExist 中）的记录
	out := make([]ScanRecord, 0, len(batch))
	for _, r := range batch {
		if _, ok := preExist[r.FullFileName]; !ok {
			out = append(out, r)
		}
	}
	return out
}

// existing 返回 fullNames 中已存在于库内的路径集合（用于筛选新增记录）。
func (m *ScanFile) existing(fullNames []string) (map[string]struct{}, error) {
	var rows []ScanRecord
	err := m.db.Where("full_file_name IN ?", fullNames).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	var set = make(map[string]struct{}, len(rows))
	for _, r := range rows {
		set[r.FullFileName] = struct{}{}
	}
	return set, nil
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

// Search 模糊查找（等价于原 SearchFromRedisSet）。
// 对 key 做同样归一化后，直接下推 SQL `LIKE` 查询 norm_file_name，
// 避免全表载入内存再过滤。返回 "path,size" 字符串。
func (m *ScanFile) Search(keyword string) []string {
	// 归一化搜索词（小写、去 -/_ ），与 norm_file_name 列同源
	kw := normalizeName(keyword)
	if kw == "" {
		return nil
	}

	// LIKE 通配符 % _ 需转义，保证子串字面匹配；SQLite 对 ASCII 的 LIKE 默认不区分大小写，
	// 但 norm 列已小写，这里以小写关键词匹配即正确。
	escaped := escapeLike(kw)
	var rows []ScanRecord
	m.db.Where("norm_file_name LIKE ? ESCAPE '\\'", "%"+escaped+"%").Order("full_file_name ASC").Find(&rows)

	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.FullFileName+","+strconv.FormatInt(r.Size, 10))
	}
	return out
}

// escapeLike 转义 LIKE 通配符 `%` 与 `_`（以及转义符本身），使关键词按字面匹配。
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
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
