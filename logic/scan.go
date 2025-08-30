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

func Save(cacheKey, filePath string, filterType map[string]struct{}) (err error) {
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
		if cache.SAdd(cacheKey, file) == 1 {
			_, _ = bw.WriteString(file + help.FormatFileSize1000(info.Size()) + "\n")
			fmt.Printf("new file %s %s\n", path, color.HiGreen.Sprint(help.FormatFileSize1000(info.Size())))
		}
		return nil
	})
	_ = bw.Flush()
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
