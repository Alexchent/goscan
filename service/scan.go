package service

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	mconf "github.com/Alexchent/goscan/config"
	"log"
	"os"
	"strings"
	"sync"
)

const CacheKey = "have_save_file"

// fd 小写不对外
var fd *os.File

func init() {
	var once sync.Once
	once.Do(func() {
		filename := "have_save_file.txt"
		//fmt.Println("打开文件\t" + filename)
		fd, _ = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	})
}

func WriteToFile(filePath string) {
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("正在扫描：", filePath)

	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		// 如果是影藏文件，直接跳过
		if fileName[0] == '.' {
			continue
		}
		if fileInfoList[i].IsDir() {
			WriteToFile(filePath + "/" + fileName)
		} else {
			// 判断是否是忽略的文件类型
			ignore := false
			for _, v := range mconf.Conf.FilterType {
				if strings.HasSuffix(fileName, v) {
					ignore = true
				}
			}
			if ignore {
				continue
			}

			filename := filePath + "/" + fileName
			// 保存到redis成功，说明是新的文件
			if mredis.SAdd(CacheKey, filename) == 1 {
				fmt.Println("发现新的文件：", filename)
				// myFile.AppendContent("have_save_file.txt", filename)
				_, err := fd.WriteString(filename + "\n")
				if err != nil {
					panic(filename + "\t文件写失败")
					return
				}
			}
		}
	}
}
