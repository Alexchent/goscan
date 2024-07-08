package help

import (
	"fmt"
	"os"
	"path/filepath"
)

// 后缀
var suffix map[string]any

func init() {
	//suffixs = []string{".git", ".DS_Store", ".localized", ".gitignore", ".sample"}
	suffix = map[string]any{
		"git": nil,
		"jpg": nil,
	}
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
			_, ok := suffix[ext]
			if ok {
				fmt.Printf("忽略 %s\n", path)
				return err
			}
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
