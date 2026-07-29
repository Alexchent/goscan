package duplicate

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Alexchent/goscan/pkg"
	"github.com/fatih/color"
)

const consumer = 100

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
				MD5:          pkg.GetFileMd5(fullName),
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
			var fileInfos []FileInfo
			var size int64
			for _, f := range files {
				fileInfos = append(fileInfos, FileInfo{
					Name:          filepath.Base(f.FullFileName),
					Path:          f.FullFileName,
					Size:          f.Size,
					SizeFormatted: pkg.FormatFileSize1000(f.Size),
				})
				size = f.Size
			}

			results = append(results, DuplicateResult{
				MD5:           md5,
				Files:         fileInfos,
				Size:          size,
				SizeFormatted: pkg.FormatFileSize1000(size),
				Count:         len(files),
				Saved:         size * int64(len(files)-1),
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
	var duplicateCount int
	for _, r := range results {
		totalSaved += r.Saved
		duplicateCount += r.Count - 1
	}
	fmt.Printf("重复文件数: %s\n", color.YellowString("%d", duplicateCount))
	fmt.Printf("可节省空间: %s\n", color.RedString(pkg.FormatFileSize1000(totalSaved)))

	fmt.Println(color.CyanString("\n----------------------------------------"))
	fmt.Println(color.CyanString("重复文件详情:"))
	fmt.Println(color.CyanString("----------------------------------------"))

	for i, r := range results {
		fmt.Printf("\n%s %s (共 %d 个副本，大小: %s)\n",
			color.BlueString("[%d]", i+1),
			color.RedString("MD5: %s...", r.MD5[:16]),
			r.Count,
			color.YellowString(r.SizeFormatted))

		for j, file := range r.Files {
			if j == 0 {
				fmt.Printf("  %s %s\n", color.GreenString("★ 原文件:"), file.Path)
			} else {
				fmt.Printf("  %s %s\n", color.RedString("  副本:"), file.Path)
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
		for j, fileInfo := range r.Files {
			row := []string{
				fmt.Sprintf("%d", i+1),
				r.MD5,
				r.SizeFormatted,
				fmt.Sprintf("%d", r.Count),
				pkg.FormatFileSize1000(r.Saved),
				fileInfo.Path,
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
