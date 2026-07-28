package service

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/Alexchent/goscan/help"
)

// FileList map的key是文件的md5，value是文件的路径
type FileList map[string][]string

func RemoveSameFile(file FileList) {
	if len(file) == 0 {
		fmt.Println("没有重复文件")
		return
	}
	for _, v := range file {
		if len(v) > 1 {
			// 保留第一个文件，删除其他文件
			for _, v := range v[1:] {
				encoded := help.MD5(v)
				cmd := exec.Command("mv", v, "/Users/chentao/same/"+encoded+"-"+path.Base(v))
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("combined out:\n%s\n", string(output))
				}
			}
		}
	}
}
