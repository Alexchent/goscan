package logic

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Alexchent/goscan/cache"
	"github.com/Alexchent/goscan/help"
	"github.com/gookit/color"
)

const CacheKey = "have_save_file"

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

func Save(cacheKey, filePath string, filterType map[string]struct{}) (err error) {
	if err = openLogfile(); err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("访问路径 %s 出错: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
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

		// 保存到redis成功，说明是新的文件
		file := fmt.Sprintf("%s,%d", path, info.Size())
		if cache.SAdd(cacheKey, file) == 1 {
			// 使用互斥锁保护写入操作
			logger.mu.Lock()
			_, _ = logger.bw.WriteString(file + help.FormatFileSize1000(info.Size()) + "\n")
			logger.mu.Unlock()
			fmt.Printf("new file %s %s\n", path, color.HiGreen.Sprint(help.FormatFileSize1000(info.Size())))
		}
		return nil
	})

	// 确保日志刷新
	logger.mu.Lock()
	_ = logger.bw.Flush()
	logger.mu.Unlock()

	return err
}

func SearchFromRedisSet(key, path string) (count int) {
	res := cache.SMembers(key)
	count = 0
	//过滤掉特殊字符-和_
	reg, _ := regexp.Compile("-|_")
	path = reg.ReplaceAllString(path, "")
	for _, val := range res {
		a := reg.ReplaceAllString(val, "")
		if strings.Contains(strings.ToLower(a), strings.ToLower(path)) {
			res := strings.Split(val, ",")
			if len(res) == 2 {
				fileSize, err := strconv.ParseInt(res[1], 10, 64)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(res[0], color.HiGreen.Sprint(help.FormatFileSize1000(fileSize)))
				}
			} else {
				fmt.Println(val)
			}
			count++
		}
	}
	return
}
