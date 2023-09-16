/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
	"os"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "加载文件中的记录",
	Long:  `导入文件中的记录`,
	Run: func(cmd *cobra.Command, args []string) {
		//rootCmd.PersistentFlags().StringVar(&readFile, "f", "", "导入的日志文件")
		//myFile.ReadString(loadFile)
		f, err := os.Open(loadFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		fileScanner := bufio.NewScanner(f)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			//fmt.Println(fileScanner.Text())
			if mredis.SAdd(scan.CacheKey, fileScanner.Text()) == 1 {
				fmt.Println("写入：", fileScanner.Text())
			}
		}
		fmt.Println("finish")
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
