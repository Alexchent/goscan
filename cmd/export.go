/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	mconf "github.com/Alexchent/goscan/config"
	myFile "github.com/Alexchent/goscan/file"
	"github.com/Alexchent/goscan/model"
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
			dir, _ := os.UserHomeDir()
			saveDir = dir + SaveDir
		}
		fmt.Println("导出文件的路径:", saveDir)

		myFile.CreateDateDir(saveDir)

		var data []string
		filename := fmt.Sprintf(saveDir+"/"+SavePath, time.Now().Unix())
		fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		data = mredis.SMembers(CacheKey)
		// 过滤掉换行符
		for _, v := range data {
			filename = strings.Trim(v, "\n")
			_, err := fd.WriteString(filename + "\n")
			if err != nil {
				panic("文件写失败：" + v)
			}

			var fileSize int64
			stat, err := os.Stat(filename)
			if err != nil {
				fileSize = 0
			} else {
				fileSize = stat.Size()
			}
			model.NewScanFile().Insert(filename, fileSize)
		}
	},
}

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "导出文件路径")
	rootCmd.AddCommand(exportCmd)
}
