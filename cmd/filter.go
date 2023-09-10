/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/redis"
	"strings"

	"github.com/spf13/cobra"
)

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "清理掉制定的文件类型",
	Long:  `清理掉制定的文件类型`,
	Run: func(cmd *cobra.Command, args []string) {
		key := "have_save_file"
		val := redis.SMembers(key)
		for _, v := range val {
			newv := strings.TrimRight(v, "\n")
			redis.SRem(key, v)
			redis.SAdd(key, newv)

			// 按 后缀清理
			if strings.HasSuffix(v, "js") ||
				strings.HasSuffix(v, "torrent") ||
				strings.HasSuffix(v, "jpeg") {
				fmt.Println("过滤掉:", v)
				redis.SRem(key, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
