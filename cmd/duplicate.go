package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Alexchent/goscan/help"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	consumer = 100
)

// DuplicateFile 重复文件信息结构
type DuplicateFile struct {
	FullFileName string
	MD5          string
	Size         int64
}

// DuplicateResult 重复文件结果
type DuplicateResult struct {
	MD5   string
	Files []string
	Size  int64
	Count int
	Saved int64 // 可节省的空间
}

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "查找重复文件",
	Long:  `查找重复文件，支持并发扫描，结果输出到控制台和CSV文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		// 检查目录参数
		if path == "" {
			fmt.Printf("请输入要扫描的目录: ")
			_, _ = fmt.Scan(&path)
		}

		_, err := os.Stat(path)
		if err != nil {
			fmt.Println(color.RedString("错误: %v", err))
			return
		}

		fmt.Println(color.CyanString("========================================"))
		fmt.Println(color.CyanString("开始扫描目录: %s", path))
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
			walkFiles(path, &scanWg, fileChan, &fileCount, &countMutex)
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

		// 输出到控制台
		printToConsole(results, fileCount, start)

		// 保存到 CSV 文件
		saveToCSV(results, path)

		fmt.Println(color.GreenString("\n✓ 扫描完成，耗时: %s", time.Since(start)))
	},
}

// walkFiles 递归遍历目录
func walkFiles(dir string, wg *sync.WaitGroup, ch chan *DuplicateFile, fileCount *int64, countMutex *sync.Mutex) {
	defer wg.Done()

	readDir, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var nextWg sync.WaitGroup
	for _, file := range readDir {
		fileName := file.Name()
		// 过滤隐藏文件
		if fileName[0] == '.' {
			continue
		}

		fullName := filepath.Join(dir, fileName)

		if file.IsDir() {
			nextWg.Add(1)
			go walkFiles(fullName, &nextWg, ch, fileCount, countMutex)
		} else {
			// 获取文件信息
			info, err := file.Info()
			if err != nil {
				continue
			}

			// 计数
			countMutex.Lock()
			*fileCount++
			countMutex.Unlock()

			// 发送到 channel
			ch <- &DuplicateFile{
				FullFileName: fullName,
				MD5:          help.GetFileMd5(fullName),
				Size:         info.Size(),
			}
		}
	}

	nextWg.Wait()
}

// processResults 处理扫描结果，筛选重复文件
func processResults(fileList map[string][]*DuplicateFile) []DuplicateResult {
	var results []DuplicateResult

	for md5, files := range fileList {
		if len(files) > 1 {
			var fileNames []string
			var size int64
			for _, f := range files {
				fileNames = append(fileNames, f.FullFileName)
				size = f.Size
			}

			results = append(results, DuplicateResult{
				MD5:   md5,
				Files: fileNames,
				Size:  size,
				Count: len(files),
				Saved: size * int64(len(files)-1), // 可节省的空间
			})
		}
	}

	return results
}

// printToConsole 输出到控制台
func printToConsole(results []DuplicateResult, totalFiles int64, start time.Time) {
	fmt.Println(color.CyanString("\n========================================"))
	fmt.Println(color.CyanString("扫描结果"))
	fmt.Println(color.CyanString("========================================"))

	fmt.Printf("总文件数: %s\n", color.YellowString("%d", totalFiles))
	fmt.Printf("重复文件组: %s\n", color.YellowString("%d", len(results)))

	if len(results) == 0 {
		fmt.Println(color.GreenString("\n✓ 没有发现重复文件"))
		return
	}

	// 计算总可节省空间
	var totalSaved int64
	for _, r := range results {
		totalSaved += r.Saved
	}
	fmt.Printf("可节省空间: %s\n", color.RedString(help.FormatFileSize1000(totalSaved)))

	fmt.Println(color.CyanString("\n----------------------------------------"))
	fmt.Println(color.CyanString("重复文件详情:"))
	fmt.Println(color.CyanString("----------------------------------------"))

	for i, r := range results {
		fmt.Printf("\n%s %s (共 %d 个副本，大小: %s)\n",
			color.BlueString("[%d]", i+1),
			color.RedString("MD5: %s...", r.MD5[:16]),
			r.Count,
			color.YellowString(help.FormatFileSize1000(r.Size)))

		for j, file := range r.Files {
			if j == 0 {
				fmt.Printf("  %s %s\n", color.GreenString("★ 原文件:"), file)
			} else {
				fmt.Printf("  %s %s\n", color.RedString("  副本:"), file)
			}
		}
	}
}

// saveToCSV 保存结果到 CSV 文件
func saveToCSV(results []DuplicateResult, scanPath string) {
	if len(results) == 0 {
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("duplicate_files_%s.csv", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(color.RedString("创建CSV文件失败: %v", err))
		return
	}
	defer file.Close()

	// 写入 UTF-8 BOM，确保 Excel 正确识别编码
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	header := []string{"序号", "MD5", "文件大小", "副本数量", "可节省空间", "文件路径"}
	if err := writer.Write(header); err != nil {
		fmt.Println(color.RedString("写入CSV表头失败: %v", err))
		return
	}

	// 写入数据
	for i, r := range results {
		for j, filePath := range r.Files {
			row := []string{
				fmt.Sprintf("%d", i+1),
				r.MD5,
				help.FormatFileSize1000(r.Size),
				fmt.Sprintf("%d", r.Count),
				help.FormatFileSize1000(r.Saved),
				filePath,
			}
			if j == 0 {
				row[0] = fmt.Sprintf("%d (原文件)", i+1)
			} else {
				row[0] = fmt.Sprintf("%d (副本%d)", i+1, j)
			}
			if err := writer.Write(row); err != nil {
				fmt.Println(color.RedString("写入CSV数据失败: %v", err))
				return
			}
		}
	}

	fmt.Printf("\n%s %s\n", color.GreenString("✓ 结果已保存到:"), color.YellowString(filename))
}

func init() {
	duplicateCmd.Flags().StringVarP(&path, "dir", "d", "", "扫描的目录")
	rootCmd.AddCommand(duplicateCmd)
}
