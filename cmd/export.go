/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	mconf "github.com/Alexchent/goscan/config"
	myFile "github.com/Alexchent/goscan/file"
	scan "github.com/Alexchent/goscan/service"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var output string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "将所有记录到处到文件中",
	Long:  `将所有记录到处到文件中`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("文件导出完成，用时：", time.Since(start))
		var saveDir string

		if output != "" {
			saveDir = output
		} else if mconf.Conf.Dir != "" {
			saveDir = mconf.Conf.Dir
		} else {
			dir, err := os.UserHomeDir()
			if err != nil {
				return
			}
			saveDir = dir + scan.SaveDir
		}
		fmt.Println("导出文件的路径:", saveDir)

		myFile.CreateDateDir(saveDir)

		var data []string
		filename := fmt.Sprintf(saveDir+scan.SavePath, time.Now().Unix())

		data = mredis.SMembers(scan.CacheKey)
		// 过滤掉换行符
		for _, v := range data {
			myFile.AppendContent(filename, strings.Trim(v, "\n"))
		}
	},
}

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "导出文件路径")
	rootCmd.AddCommand(exportCmd)
}
