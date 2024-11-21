package cmd

import (
	_ "embed"
	"github.com/Alexchent/goscan/service"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var outputConf string

// cleanCmd represents the clean command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "生成默认配置文件",
	Long:  `生成默认配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		service.MakeConfig(outputConf)
	},
}

// go run main clean -c apk
func init() {
	// 本地标志, 此处进队cleanCmd有效
	home, _ := homedir.Dir()
	configCmd.Flags().StringVarP(&outputConf, "output", "o", home, "配置文件输出路径")
	rootCmd.AddCommand(configCmd)
}
