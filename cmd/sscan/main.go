// Command sscan 是可 go install 安装的 sscan 入口。
//
//	go install github.com/Alexchent/goscan/cmd/sscan@latest
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/service"
)

func main() {
	if err := config.InitConf("", ""); err != nil {
		fmt.Println(err)
		return
	}
	service.SetDB(config.Conf.DB)

	// 关键字取自命令行参数；缺失时交互式读取
	keyword := ""
	if len(os.Args) > 1 {
		keyword = strings.Join(os.Args[1:], " ")
	}
	if len(keyword) == 0 {
		fmt.Printf("请输入要查询的文件:\n")
		if _, err := fmt.Scan(&keyword); err != nil {
			return
		}
	}
	count := service.Search(keyword)
	fmt.Printf("本次扫描发现 %d 个文件\n", count)
}
