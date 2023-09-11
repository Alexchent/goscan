package myFile

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func AppendContent(filename, content string) {
	//filename := fmt.Sprintf(SavePath, time.Now().Format("060102"))
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 打开文件失败处理
		log.Fatal(err.Error())
	} else {
		//buf := []byte(content)
		//fd.Write(buf)
		fd.WriteString(content + "\n")
	}

}

// CreateDateDir basePath是固定目录路径
func CreateDateDir(folderPath string) (dirPath string) {
	//folderName := time.Now().Format("2006-01-02")
	//folderPath := filepath.Join(basePath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步
		// 先创建文件夹
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		// 再修改权限
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			log.Fatal(err)
			return ""
		}
	}
	return folderPath
}

func readString(filename string) {
	//f, err := os.Open("/Users/chentao/Downloads/down.txt")
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && io.EOF == err {
			break
		}
		data := strings.Trim(line, "\n")

		if data != "" {
			fmt.Println(data)
			// 写入导redis内
			//redis.SAdd(service.CacheKey, data)
		}
	}
}
