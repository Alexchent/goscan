package main

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/redis"
	"strings"
)

func main() {
	var path string
	fmt.Printf("请输入要查询的文件:\n")
	_, err := fmt.Scan(&path)
	if err != nil {
		return
	}
	path = strings.ToLower(path)
	count := 0
	count += SearchFromRedisSet("have_save_file", path)
	count += SearchFromRedisSet("laravel_database_files", path)
	res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
	fmt.Println(res)
}

func SearchFromRedisSet(key, path string) (count int) {
	res := redis.SMembers(key)
	count = 0
	for _, val := range res {
		if strings.Contains(strings.ToLower(val), path) {
			fmt.Println(val)
			count++
		}
	}
	return
}
