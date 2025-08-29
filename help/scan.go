package help

import (
	"fmt"
	"os"
	"path/filepath"
)

// 后缀
var suffix map[string]any

func init() {
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
			if info.Name()[0] == '.' {
				return err
			}
			ext := filepath.Ext(path)
			_, ok := suffix[ext]
			if ok {
				return err
			}
			//fmt.Println(info.Name(), info.Size(), info.ModTime().Format(time.DateTime))
			fmt.Printf("文件大小: %s, 路径: %s\n", FormatFileSize(info.Size()), path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
