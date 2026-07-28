package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/Alexchent/goscan/help"
)

//go:embed example.yaml
var configData string

func MakeConfig(outputConf string) {
	help.CreateDateDir(outputConf)
	filename := outputConf + "/.scan.yaml"

	if help.FileIsExist(filename) {
		fmt.Println("配置文件已存在")
		return
	}

	fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	_, err := fd.WriteString(configData)
	if err != nil {
		fmt.Println("配置文件生成失败")
	}
}
