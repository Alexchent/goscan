package logic

import (
	"bufio"
	"fmt"
	"github.com/Alexchent/goscan/cache"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const CacheKey = "have_save_file"

// fd 小写不对外
// var fd *os.File
var bw *bufio.Writer
var once sync.Once

func openLogfile() {
	once.Do(func() {
		filename := "have_save_file.log"
		fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		bw = bufio.NewWriter(fd)
	})
}

func Save(filePath string, filterType map[string]struct{}) (err error) {
	openLogfile()
	err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
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
		if cache.SAdd(CacheKey, file) == 1 {
			_, _ = bw.WriteString(file + "\n")
		} else {
			fmt.Println("file is exist ", path)
		}
		return nil
	})
	_ = bw.Flush()
	return err
}
