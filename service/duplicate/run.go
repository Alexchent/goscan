package duplicate

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Run 查找重复文件的入口：校验目录、并发扫描、收集结果、输出控制台、
// 有重复时启动 Web 服务器。行为与原 cmd/duplicate.go 的 Run 一致。
func Run(scanDir string) {
	start := time.Now()

	// 检查目录参数
	if scanDir == "" {
		fmt.Printf("请输入要扫描的目录: ")
		_, _ = fmt.Scan(&scanDir)
	}

	_, err := os.Stat(scanDir)
	if err != nil {
		fmt.Println(color.RedString("错误: %v", err))
		return
	}

	fmt.Println(color.CyanString("========================================"))
	fmt.Println(color.CyanString("开始扫描目录: %s", scanDir))
	fmt.Println(color.CyanString("========================================"))

	// 使用 channel 进行并发处理
	fileChan := make(chan *DuplicateFile, consumer)
	var scanWg sync.WaitGroup
	scanWg.Add(1)

	// 文件计数器
	var fileCount int64
	var countMutex sync.Mutex

	// 启动扫描 goroutine
	go func() {
		walkFiles(scanDir, &scanWg, fileChan, &fileCount, &countMutex)
	}()

	// 收集结果
	fileList := make(map[string][]*DuplicateFile)
	var collectWg sync.WaitGroup
	collectWg.Add(1)

	go func() {
		defer collectWg.Done()
		for f := range fileChan {
			fileList[f.MD5] = append(fileList[f.MD5], f)
		}
	}()

	// 等待扫描完成
	scanWg.Wait()
	close(fileChan)
	collectWg.Wait()

	// 处理结果
	results := processResults(fileList)
	totalFilesScanned = fileCount
	scanResults = results

	// 输出到控制台
	printToConsole(results, fileCount, start)

	// 保存到 CSV 文件
	// saveToCSV(results, scanDir)

	// 只有找到重复文件时才启动 Web 服务器
	if len(results) > 0 {
		startWebServer()
	}
}
