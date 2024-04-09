package cmd

import (
	"github.com/Alexchent/goscan/config"
	"github.com/spf13/cobra"
)

const SaveDir = "/scanLog"
const SavePath = "have_save_file_%d.txt"
const CacheKey = "have_save_file"
const CacheKeyMd5 = "have_save_file_md5"

var (
	// Used for flags.
	cfgFile string
	path    string

	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "scan 用于扫描文件目录，和查找文件",
		Long:  `scan 是一个用于扫描本地文件，形成日志文件，以便或许快速查找文件的程序`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// 绑定参数示例
// go run main.go -c=./scan.yaml
func init() {

	// 标志可以是 "persistent" 的，这意味着该标志将可用于分配给它的命令以及该命令下的每个命令。对于全局标志，将标志分配为根上的持久标志。
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config.InitConf(cfgFile)
}
