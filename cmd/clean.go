package cmd

import (
	"fmt"

	"github.com/Alexchent/goscan/service"

	"github.com/spf13/cobra"
)

var clearSuffix string
var clearKey string

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理掉符合条件的的文件",
	Long:  `清理掉符合条件的的文件`,
	Run: func(cmd *cobra.Command, args []string) {
		// 按后缀清理
		if clearSuffix != "" {
			for _, v := range service.DeleteBySuffix(clearSuffix) {
				fmt.Println("过滤掉：", v)
			}
		}

		// 按关键字清理
		if clearKey != "" {
			for _, v := range service.DeleteByKeyword(clearKey) {
				fmt.Println("过滤掉：", v)
			}
		}
	},
}

// 后缀
// go run main clean -c apk
func init() {
	// 本地标志, 此处进队cleanCmd有效
	cleanCmd.Flags().StringVar(&clearSuffix, "suffix", "", "需要清理掉的文件类型")
	cleanCmd.Flags().StringVarP(&clearKey, "keywords", "k", "", "关键字")
	rootCmd.AddCommand(cleanCmd)
}
