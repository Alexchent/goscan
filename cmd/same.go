package cmd

import (
	"fmt"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"time"
)

// startCmd represents the start command
var sameCmd = &cobra.Command{
	Use:   "same",
	Short: "查找重复文件",
	Long:  `查找重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))

		if path == "" {
			fmt.Printf("请输入要扫描的目录:\n")
			_, _ = fmt.Scan(&path)
		}

		fmt.Println("开始扫描：" + path)
		scan.GroupByFileMD5(path)
		// 输出重复文件到txt中
		fmt.Println("记录扫描文件 start")
		same := scan.LogSameFile()
		fmt.Println("记录扫描文件 end")
		if same == nil {
			return
		}
		//fmt.Println("执行删除操作 start")
		//scan.RemoveSameFile(same)
		//fmt.Println("执行删除操作 end")
	},
}

func init() {
	sameCmd.Flags().StringVarP(&path, "dir", "d", "", "需要扫描的目录")
	rootCmd.AddCommand(sameCmd)
}
