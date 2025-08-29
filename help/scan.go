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
			ext := filepath.Ext(path)
			_, ok := suffix[ext]
			if ok {
				//fmt.Printf("忽略 %s\n", path)
				return err
			}
			//fmt.Println(info.Name(), info.Size(), info.ModTime().Format(time.DateTime))
			//return nil
			// 根据文件大小自动转化为MB或GB
			if info.Size() > 1024*1024*1024 {
				fmt.Printf("文件大小: %.2f GB, 路径: %s\n", float64(info.Size())/1024/1024/1024, path)
				return nil
			}
			if info.Size() > 1024*1024 {
				fmt.Printf("文件大小: %.2f MB, 路径: %s\n", float64(info.Size())/1024/1024, path)
				return nil
			}
			if info.Size() > 1024 {
				fmt.Printf("文件大小: %d KB, 路径: %s\n", info.Size()/1024, path)
				return nil
			}
			fmt.Printf("文件大小: %d B, 路径: %s\n", info.Size(), path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
