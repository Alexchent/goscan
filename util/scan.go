package help

import (
	"fmt"
	"os"
	"path/filepath"
)

var suffixs []string

func init() {
	suffixs = []string{".git", ".DS_Store", ".localized", ".gitignore", ".sample"}
}

func Scale(dir string) error {
	//获取当前目录下的所有文件或目录信息
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			//fmt.Println(path)
			ext := filepath.Ext(path)
			for _, suffix := range suffixs {
				if ext == suffix {
					fmt.Printf("忽略 %s\n", path)
					continue
				}
			}
			appendContent(path)

		}

		return nil
	})

	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

func appendContent(content string) {
	fd, err := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 打开文件失败处理
	} else {
		content += "\n"
		buf := []byte(content)
		_, err := fd.Write(buf)
		if err != nil {
			return
		}
	}
}
