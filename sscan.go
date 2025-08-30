package main

import (
	"fmt"
	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/logic"
	"github.com/mitchellh/go-homedir"
)

func main() {
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cfgFile := dir + "/.scan.yaml"
	config.InitConf(cfgFile)

	var path string
	fmt.Printf("请输入要查询的文件:\n")
	_, err = fmt.Scan(&path)
	if err != nil {
		return
	}
	count := 0
	count += logic.SearchFromRedisSet(logic.CacheKey, path)
	res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
	fmt.Println(res)
}
