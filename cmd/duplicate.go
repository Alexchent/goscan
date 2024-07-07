package cmd

import (
	"github.com/Alexchent/goscan/cmd/duplicate"
	"github.com/spf13/cobra"
)

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "查找重复文件",
	Long:  `查找重复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		duplicate.Do(path)
	},
}

func init() {
	duplicateCmd.Flags().StringVarP(&path, "dir", "d", "", "扫描的目录")
	rootCmd.AddCommand(duplicateCmd)
}
