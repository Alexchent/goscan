package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexchent/goscan/logic"

	"github.com/spf13/cobra"
)

const (
	saveFileFmt = "have_save_file_%d.txt"
)

var output string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "dump",
	Short: "导出历史扫描记录到本地文件",
	Long: `将历史扫描记录导出到本地文本文件，每条记录占一行。

默认导出到当前目录，文件名格式为 have_save_file_<Unix时间戳>.txt。
可通过 -o/--output 指定导出目录（路径必须以路径分隔符结尾，如 ./logs/ 或 /tmp/），
否则仍导出到当前目录。`,
	Example: `  # 导出到当前目录
  scan dump

  # 导出到指定目录（注意结尾的路径分隔符）
  scan dump -o ./logs/
  scan dump --output /tmp/`,
	Run: func(cmd *cobra.Command, args []string) {
		saveDir := "."
		if len(output) > 0 && os.IsPathSeparator(output[len(output)-1]) {
			saveDir = output
		}

		// 导出文件的全路径
		exportPath := fmt.Sprintf(saveDir+"/"+saveFileFmt, time.Now().Unix())
		fd, err := os.OpenFile(exportPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("打开导出文件失败: %v\n", err)
			return
		}
		defer func() {
			_ = fd.Close()
		}()

		start := time.Now()
		defer func() {
			fmt.Println("文件导出完成，用时：", time.Since(start))
		}()

		data := logic.AllRecords()
		// 过滤掉换行符
		for _, v := range data {
			line := strings.Trim(v, "\n")
			if _, err := fd.WriteString(line + "\n"); err != nil {
				fmt.Printf("文件写失败: %v (内容: %s)\n", err, line)
				return
			}
		}
	},
}

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "导出目录路径（须以路径分隔符结尾，默认为当前目录）")
	rootCmd.AddCommand(exportCmd)
}
