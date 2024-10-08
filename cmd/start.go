package cmd

import (
	"fmt"
	mconf "github.com/Alexchent/goscan/config"
	scan "github.com/Alexchent/goscan/service"
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "开始扫描",
	Long:  `开始扫描`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer fmt.Println("扫描完成，耗时：", time.Since(start))
		fmt.Println("ignore:", mconf.Conf.FilterType)
		if path == "" {
			fmt.Printf("请输入要扫描的目录:\n")
			_, err := fmt.Scan(&path)
			if err != nil {
				return
			}
			if path == "/" || path == "." {
				dir, err := homedir.Dir()
				if err != nil {
					panic(err)
				}
				path = dir + "/Downloads"
			}
		}
		path = strings.TrimRight(path, "/")
		color.HiGreen.Println("开始扫描：", path)
		scan.WriteToFile(path)
	},
}

func init() {
	startCmd.Flags().StringVarP(&path, "dir", "d", "", "需要扫描的目录")
	rootCmd.AddCommand(startCmd)
}
