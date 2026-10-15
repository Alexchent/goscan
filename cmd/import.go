package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var (
	loadFile string

	loadCmd = &cobra.Command{
		Use:   "import",
		Short: "加载文件中的记录",
		Long:  `导入文件中的记录`,
		Run: func(cmd *cobra.Command, args []string) {

			f, err := os.Open(loadFile)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer f.Close()

			fileScanner := bufio.NewScanner(f)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				if fileScanner.Text() != "" {
					if service.LoadRecord(fileScanner.Text()) {
						fmt.Println("写入：", fileScanner.Text())
					}
				}
			}
			if err = fileScanner.Err(); err != nil {
				fmt.Println("scan file err: " + err.Error())
			}
			fmt.Println("finish")
		},
	}
)

func init() {
	loadCmd.Flags().StringVarP(&loadFile, "loadFile", "f", "", "导入的日志文件")
	rootCmd.AddCommand(loadCmd)
}
