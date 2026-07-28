package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/help"
	"github.com/Alexchent/goscan/logic"

	"github.com/spf13/cobra"
)

// 这两个常量仅在 dump 命令中使用，从原 root.go 迁移至此。
// saveDirSuffix 只是路径后缀片段，必须拼接到某个 base 目录之后使用，
// 切勿当作完整目录路径直接传给 os.MkdirAll 之类，否则会操作根目录 "/scanLog"。
const (
	saveDirSuffix = "/scanLog"
	saveFileFmt   = "have_save_file_%d.txt"
)

var output string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "dump",
	Short: "导出记录",
	Long:  `将历史记录导出到文件中`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			fmt.Println("文件导出完成，用时：", time.Since(start))
		}()

		var saveDir string

		if output != "" {
			saveDir = output
		} else if config.Conf.Dir != "" {
			saveDir = config.Conf.Dir
		} else {
			dir, _ := os.UserHomeDir()
			saveDir = dir + saveDirSuffix
		}
		fmt.Println("导出文件的路径:", saveDir)

		help.CreateDateDir(saveDir)

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

		data := logic.AllRecords()
		//svc := logic.NewSaveLogic(*mconf.Conf)
		// 过滤掉换行符
		for _, v := range data {
			line := strings.Trim(v, "\n")
			if _, err := fd.WriteString(line + "\n"); err != nil {
				fmt.Printf("文件写失败: %v (内容: %s)\n", err, line)
				return
			}
			//svc.Save(exportPath)
		}
	},
}

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "导出文件路径")
	rootCmd.AddCommand(exportCmd)
}
