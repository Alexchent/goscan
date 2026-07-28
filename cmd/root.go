package cmd

import (
	"fmt"
	"os"

	"github.com/Alexchent/goscan/config"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string
	// scanDir 是 --dir flag 的接收变量，仅在 start / duplicate 子命令中使用，
	// 不应与标准库 path 包或局部变量混用，请勿在其它文件中再声明同名局部变量。
	scanDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan 用于扫描文件目录，和查找文件",
	Long:  `scan 是一个用于扫描本地文件，形成日志文件，以便或许快速查找文件的程序`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// 绑定参数示例
// go run main.go -c=.scan.yaml
func init() {
	// 标志可以是 "persistent" 的，这意味着该标志将可用于分配给它的命令以及该命令下的每个命令。对于全局标志，将标志分配为根上的持久标志。
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if err := config.InitConf(cfgFile); err != nil {
		// 在 OnInitialize 回调里直接 os.Exit 会跳过所有 defer，
		// 这里保持原有行为，但建议后续改为通过 cobra 的 RunE 返回 error。
		fmt.Printf("配置初始化失败: %v\n", err)
		os.Exit(1)
	}
}
