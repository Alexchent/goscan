package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	loadFile    string

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
// go run main.go version --license=MIT 等同于 go run main.go version -l=MIT
func init() {
	cobra.OnInitialize(initConfig)

	// 标志可以是 "persistent" 的，这意味着该标志将可用于分配给它的命令以及该命令下的每个命令。对于全局标志，将标志分配为根上的持久标志。
	// --config="/config/scan/.cobra.yaml"
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	// 绑定到变量 userLicense
	// --license=MIT 或 -l=license
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "许可证")
	rootCmd.PersistentFlags().StringVarP(&loadFile, "loadFile", "f", "", "导入的日志文件")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
