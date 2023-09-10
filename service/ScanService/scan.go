package scan

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	myFile "github.com/Alexchent/goscan/file"
	"log"
	"os"
	"os/user"
	path "path/filepath"
)

//const SavePath = "/Users/chentao/have_save_file_%s.txt"

const SaveDir = "/scanLog/"
const SavePath = "have_save_file_%d.txt"
const CacheKey = "have_save_file"

type Config struct {
	IgnoreFileType []string // 忽略的文件类型
	SaveDir        string   // 保存到文件目录
}

var config Config

func init() {
	// 从toml配置文件中读取配置，绑定到Config结构体
	config.IgnoreFileType = []string{".js", ".torrent", ".jpg", ".jpeg", ".localized"}
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	config.SaveDir = u.HomeDir
}

func (c *Config) GetSaveDir() string {
	return c.SaveDir
}

func WriteToFile(filePath string) {
	fileInfoList, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("正在扫描：", filePath)

	for i := range fileInfoList {
		fileName := fileInfoList[i].Name()
		if fileInfoList[i].IsDir() {
			WriteToFile(filePath + "/" + fileName)
		} else {
			// 过滤Mac的.DS_Store文件
			if fileInfoList[i].Name() == ".DS_Store" {
				continue
			}
			// 过滤js和torrent、jpg文件
			baseName := path.Base(fileName)
			ext := path.Ext(baseName)

			//if ext == ".js" || ext == ".torrent" || ext == ".jpg" || ext == ".jpeg" {
			//	continue
			//}

			// 判断是否是忽略的文件类型
			// 判断ext是否在IgnoreFileType中
			ignore := false
			for _, v := range config.IgnoreFileType {
				if ext == v {
					//log.Println("忽略文件：", fileName)
					ignore = true
				}
			}

			if ignore {
				continue
			}

			content := filePath + "/" + fileName

			// 保存到redis成功，说明是新的文件
			if mredis.SAdd(CacheKey, content) == 1 {
				fmt.Println("发现新的文件：", fileName)
				//filename := fmt.Sprintf(SavePath, time.Now().Format("060102"))
				myFile.AppendContent("have_save_file.txt", content)
			}
		}
	}
}
