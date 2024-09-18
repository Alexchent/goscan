package cmd

import (
	"encoding/json"
	"fmt"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var file string

// startCmd represents the start command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "移动重复文件",
	Long:  `移动重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))

		if file == "" {
			fmt.Println(file + "\t文件不存在")
			return
		}

		readFile, err := os.ReadFile(file)
		if err != nil {
			panic("打开文件失败: " + err.Error())
			return
		}

		var same scan.FileList
		err = json.Unmarshal(readFile, &same)
		if err != nil {
			log.Fatal("error during unmarshal():", err)
			return
		}
		//fmt.Println(same)
		fmt.Println("执行删除操作 start")
		scan.RemoveSameFile(same)
		fmt.Println("执行删除操作 end")
	},
}

func init() {
	removeCmd.Flags().StringVarP(&file, "file", "f", "", "same命令生成的文件")
	rootCmd.AddCommand(removeCmd)
}
