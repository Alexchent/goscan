package main

import (
	"fmt"
	scan "github.com/Alexchent/goscan/scan/ScanService"
	"time"
)

func main() {
	start := time.Now()
	defer fmt.Println(time.Since(start))

	var path string
	fmt.Printf("请输入要扫描的目录:\n")

	_, err := fmt.Scan(&path)
	if err != nil {
		return
	}
	if path == "" || path == "/" {
		path = "/Users/chentao/Downloads"
	}
	scan.WriteToFile(path)
}
