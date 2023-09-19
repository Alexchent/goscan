package service

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	mconf "github.com/Alexchent/goscan/config"
	myFile "github.com/Alexchent/goscan/file"
	"log"
	"os"
	"strings"
)

const SaveDir = "/scanLog/"
const SavePath = "have_save_file_%d.txt"
const CacheKey = "have_save_file"

func WriteToFile(filePath string) {
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("正在扫描：", filePath)

	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		if fileInfoList[i].IsDir() {
			WriteToFile(filePath + "/" + fileName)
		} else {
			if fileInfoList[i].Name() == ".DS_Store" {
				continue
			}

			// 判断是否是忽略的文件类型
			ignore := false
			for _, v := range mconf.Conf.FilterType {
				if strings.HasSuffix(fileName, v) {
					log.Println("忽略文件：", fileName)
					ignore = true
				}
			}
			if ignore {
				continue
			}

			content := filePath + "/" + fileName

			// 保存到redis成功，说明是新的文件
			if mredis.SAdd(CacheKey, content) == 1 {
				fmt.Println("发现新的文件：", content)
				//filename := fmt.Sprintf(SavePath, time.Now().Format("060102"))
				myFile.AppendContent("have_save_file.txt", content)
			}
		}
	}
}
