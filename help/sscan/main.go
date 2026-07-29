package main

import (
	"fmt"
	"os"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/service"
)

func main() {
	if err := config.InitConf("", ""); err != nil {
		fmt.Println(err)
		return
	}
	service.SetDB(config.Conf.DB)
	keyword := os.Args[1]
	fmt.Print(keyword)
	if len(keyword) == 0 {
		fmt.Printf("请输入要查询的文件:\n")
		_, err := fmt.Scan(&keyword)
		if err != nil {
			return
		}
	}
	count := 0
	count += service.Search(keyword)
	res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
	fmt.Println(res)
}
