/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/redis"
	myFile "github.com/Alexchent/goscan/file"
	scan "github.com/Alexchent/goscan/scan/ScanService"
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
		defer fmt.Println(time.Since(start))

		var data []string
		filename := fmt.Sprintf(scan.SavePath, time.Now().Format("060102"))

		data = redis.SMembers("have_save_file")
		for _, v := range data {
			myFile.AppendContent(filename, strings.Trim(v, "\n"))
		}

		data = redis.SMembers("laravel_database_files")
		for _, v := range data {
			myFile.AppendContent(filename, v+"\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
