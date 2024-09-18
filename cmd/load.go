package cmd

import (
	"bufio"
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	"github.com/spf13/cobra"
	"os"
)

// loadCmd represents the load command
var (
	loadFile string

	loadCmd = &cobra.Command{
		Use:   "load",
		Short: "加载文件中的记录",
		Long:  `导入文件中的记录`,
		Run: func(cmd *cobra.Command, args []string) {

			f, err := os.Open(loadFile)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			fileScanner := bufio.NewScanner(f)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				if fileScanner.Text() != "" {
					if mredis.SAdd(CacheKey, fileScanner.Text()) == 1 {
						fmt.Println("写入：", fileScanner.Text())
					}
				}
			}
			fmt.Println("finish")
		},
	}
)

func init() {
	loadCmd.Flags().StringVarP(&loadFile, "loadFile", "f", "", "导入的日志文件")
	rootCmd.AddCommand(loadCmd)
}
