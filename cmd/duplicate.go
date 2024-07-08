package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cmd/duplicate"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "查找重复文件",
	Long:  `查找重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完毕，共耗时：", time.Since(start))
		_, err := os.Stat(path)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ch := make(chan *duplicate.File)
		go duplicate.Do(path, ch)

		// 注意只有channel关闭时，才会退出
		fileList := make(duplicate.FIleList)
		for f := range ch {
			//fmt.Println(f.MD5, f.FullFileName)
			fileList[f.MD5] = append(fileList[f.MD5], f.FullFileName)
		}

		for _, filename := range fileList {
			if len(filename) > 1 {
				fmt.Println(filename)
			}
		}
		fmt.Println("over")
	},
}

func init() {
	duplicateCmd.Flags().StringVarP(&path, "dir", "d", "", "扫描的目录")
	rootCmd.AddCommand(duplicateCmd)
}
