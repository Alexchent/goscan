package service

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Alexchent/goscan/model"
	"github.com/Alexchent/goscan/pkg"
	"github.com/gookit/color"
)

// scanFile 是包级共享的 SQLite 访问入口，由 config.InitConf 通过 SetDB 注入。
// 之所以用包级变量而非每次 new：scan/find/clean/dump/load 等子命令都共享
// 同一个 scan.db，且原 Redis 实现也是全局 Client，保持同等生命周期。
var scanFile *model.ScanFile

// SetDB 注入由 config 初始化好的 ScanFile，供本包各命令复用。
func SetDB(s *model.ScanFile) {
	scanFile = s
}

// 文件日志管理器
type fileLogger struct {
	fd *os.File
	bw *bufio.Writer
	mu sync.Mutex
}

var (
	logger *fileLogger
	once   sync.Once
)

// openLogfile 打开日志文件，使用单例模式
func openLogfile() error {
	var err error
	once.Do(func() {
		filename := "have_save_file.log"
		logger = &fileLogger{}
		logger.fd, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return
		}
		logger.bw = bufio.NewWriter(logger.fd)
	})
	return err
}

// CloseLogfile 关闭日志文件，应在程序退出时调用
func CloseLogfile() error {
	if logger == nil {
		return nil
	}
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.bw != nil {
		_ = logger.bw.Flush()
	}
	if logger.fd != nil {
		return logger.fd.Close()
	}
	return nil
}

// Save 遍历 filePath，将未过滤的新文件写入 SQLite(scan_record) 与日志文件。
// filterType 为要忽略的后缀集合。返回遍历过程中的错误。
func Save(filePath string, filterType map[string]struct{}) (err error) {
	if scanFile == nil {
		return fmt.Errorf("数据库未初始化")
	}
	if err = openLogfile(); err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer CloseLogfile()

	// 批量写入缓存：攒满 batchSize 一次入库，避免扫描大量文件时逐条磁盘 IO
	const batchSize = 500
	pending := make([]*model.ScanRecord, 0, batchSize)

	flush := func() error {
		if len(pending) == 0 {
			return nil
		}
		added := scanFile.AddRecords(pending)
		for _, rec := range added {
			line := fmt.Sprintf("%s,%d", rec.FullFileName, rec.Size)
			logger.mu.Lock()
			_, _ = logger.bw.WriteString(line + pkg.FormatFileSize1000(rec.Size) + "\n")
			logger.mu.Unlock()
			fmt.Printf("new file %s %s\n", rec.FullFileName, color.HiGreen.Sprint(pkg.FormatFileSize1000(rec.Size)))
		}
		pending = pending[:0]
		return nil
	}

	err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("访问路径 %s 出错: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			// 如果是影藏目录，跳过整个目录及其子项
			if info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}
		// 如果是影藏文件，直接跳过
		if info.Name()[0] == '.' {
			return nil
		}

		// 忽略指定格式的文件
		suffix := strings.ToLower(strings.TrimPrefix(filepath.Ext(info.Name()), "."))
		if _, ok := filterType[suffix]; ok {
			return nil
		}

		// 攒入批次，攒满即入库
		pending = append(pending, &model.ScanRecord{
			FullFileName: path,
			Size:         info.Size(),
		})
		if len(pending) >= batchSize {
			return flush()
		}
		return nil
	})
	// 收尾：刷新剩余批次
	if ferr := flush(); err == nil && ferr != nil {
		err = ferr
	}

	// 确保日志刷新
	logger.mu.Lock()
	_ = logger.bw.Flush()
	logger.mu.Unlock()

	return err
}

// Search 关键字模糊查找，返回匹配条目数并打印到标准输出。
// 替换原 SearchFromRedisSet，匹配语义保持一致（去掉 -/_ 后大小写不敏感包含）。
func Search(keyword string) (count int) {
	if scanFile == nil {
		fmt.Println("数据库未初始化")
		return 0
	}
	for _, val := range scanFile.Search(keyword) {
		res := strings.Split(val, ",")
		if len(res) == 2 {
			size, err := strconv.ParseInt(res[1], 10, 64)
			if err != nil {
				fmt.Println(val)
			} else {
				fmt.Println(res[0], color.HiGreen.Sprint(pkg.FormatFileSize1000(size)))
			}
		} else {
			fmt.Println(val)
		}
		count++
	}
	return
}

// AllRecords 返回全部 "path,size" 记录，替换原 cache.SMembers。
func AllRecords() []string {
	if scanFile == nil {
		return nil
	}
	return scanFile.AllRecords()
}

// LoadRecord 从一行日志文本 "path,size"(或纯 path) 解析并入库，返回是否新写入。
// 供 load 子命令使用，替换原 cache.SAdd。
func LoadRecord(line string) bool {
	if scanFile == nil {
		return false
	}
	path, size := parseLogLine(line)
	if path == "" {
		return false
	}
	return scanFile.AddRecord(path, size)
}

// DeleteBySuffix 按后缀清理记录，返回被删除的条目列表。
func DeleteBySuffix(suffix string) []string {
	if scanFile == nil {
		return nil
	}
	return scanFile.DeleteBySuffix(suffix)
}

// DeleteByKeyword 按关键字清理记录，返回被删除的条目列表。
func DeleteByKeyword(keyword string) []string {
	if scanFile == nil {
		return nil
	}
	return scanFile.DeleteByKeyword(keyword)
}

// parseLogLine 解析 "path,size" 形式的日志行；size 缺失时返回 0。
func parseLogLine(line string) (path string, size int64) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", 0
	}
	// 兼容原始 log: 文件名本身可能不含逗号，但记录格式是 "path,size"
	// 从右侧找最后一个逗号切分，避免路径中含逗号被误切
	idx := strings.LastIndex(line, ",")
	if idx < 0 {
		return line, 0
	}
	path = line[:idx]
	// size 部分可能带日志后缀(原 log 写入时拼了 FormatFileSize1000)，
	// 仅取前导数字
	sizeStr := line[idx+1:]
	var n int64
	for _, c := range sizeStr {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int64(c-'0')
	}
	return path, n
}
