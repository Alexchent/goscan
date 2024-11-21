package logic

import (
	"bufio"
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	"github.com/Alexchent/goscan/config"
	"log"
	"os"
	"path"
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

func WriteToFile(filePath string) {
	openLogfile()

	filePath = strings.TrimRight(filePath, "/")
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Println(err)
	}
	//svc := NewSaveLogic(*mconf.Conf)
	//fmt.Println("正在扫描：", filePath)
	for _, file := range fileInfoList {
		fileName := file.Name()
		// 如果是影藏文件，直接跳过
		if fileName[0] == '.' {
			continue
		}

		if file.IsDir() {
			WriteToFile(filePath + "/" + fileName)
		} else {
			// 忽略指定格式的文件
			suffix := path.Ext(fileName)[1:]
			if _, ok := config.FilterSuffix[strings.ToLower(suffix)]; ok {
				continue
			}
			filename := filePath + "/" + fileName

			// 保存到redis成功，说明是新的文件
			if mredis.SAdd(CacheKey, filename) == 1 {
				fmt.Println(CacheKey, filename)
				bw.WriteString(filename + "\n")
				//svc.Save(filename)
			}
		}
	}
	bw.Flush()
}
