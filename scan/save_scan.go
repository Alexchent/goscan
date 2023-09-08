package main

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/redis"
	myFile "github.com/Alexchent/goscan/file"
	scan "github.com/Alexchent/goscan/scan/ScanService"
	"strings"
	"time"
)

// 备份缓存中的数据到本地
func main() {
	start := time.Now()
	defer fmt.Println(time.Since(start))

	var data []string
	filename := fmt.Sprintf(scan.SavePath, time.Now().Format("060102"))

	data = redis.SMembers("have_save_file")
	for _, v := range data {
		myFile.AppendContent(filename, strings.Trim(v, "\n"))
	}

	data = redis.SMembers("laravel_database_files")
	for _, v := range data {
		myFile.AppendContent(filename, v+"\n")
	}
}
