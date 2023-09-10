/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/redis"
	myFile "github.com/Alexchent/goscan/file"
	scan "github.com/Alexchent/goscan/service/ScanService"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "将所有记录到处到文件中",
	Long:  `将所有记录到处到文件中`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("文件导出完成，用时：", time.Since(start))

		dir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		saveDir := dir + scan.SaveDir
		myFile.CreateDateDir(saveDir)

		fmt.Println("导出文件的路径:", saveDir)
		var data []string
		filename := fmt.Sprintf(saveDir+scan.SavePath, time.Now().Unix())

		data = redis.SMembers("have_save_file")
		// 过滤掉换行符
		for _, v := range data {
			myFile.AppendContent(filename, strings.Trim(v, "\n"))
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
