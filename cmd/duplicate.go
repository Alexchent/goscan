package cmd

import (
	"github.com/Alexchent/goscan/service/duplicate"
	"github.com/spf13/cobra"
)

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "查找重复文件",
	Long:  `查找重复文件，支持并发扫描，结果输出到控制台、CSV文件和Web页面`,
	Run: func(cmd *cobra.Command, args []string) {
		duplicate.Run(scanDir)
	},
}

func init() {
	duplicateCmd.Flags().StringVarP(&scanDir, "dir", "d", "", "扫描的目录")
	rootCmd.AddCommand(duplicateCmd)
}
