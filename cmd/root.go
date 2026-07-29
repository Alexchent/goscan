package cmd

import (
	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/logic"
	"github.com/spf13/cobra"
)

var (
	// cfgFile 由 -c/--config 指定的配置文件路径（最高优先级）
	cfgFile string
	// workspace 由 -w/--workspace 指定的工作区目录，默认 ~/.scan
	workspace string
	// scanDir 是 --dir flag 的接收变量，仅在 start / duplicate 子命令中使用，
	// 不应与标准库 path 包或局部变量混用，请勿在其它文件中再声明同名局部变量。
	scanDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan 用于扫描文件目录，和查找文件",
	Long:  `scan 是一个用于扫描本地文件，形成日志文件，以便或许快速查找文件的程序`,
	// PersistentPreRunE 对所有子命令生效：先解析工作区并加载配置。
	// init 命令负责创建工作区，不要求工作区已初始化，故豁免。
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 以下命令不依赖工作区配置，豁免初始化校验：
		//   init      —— 负责创建工作区
		//   duplicate —— 纯内存扫描，不读写 DB
		//   remove    —— 读 JSON 执行 mv，不依赖工作区
		switch cmd.Name() {
		case "init", "duplicate", "remove":
			return nil
		}
		// version 允许在工作区未初始化时降级运行（只打印版本号），
		// 已初始化则正常加载以展示完整配置信息。
		if cmd.Name() == "version" {
			_ = loadConfig()
			return nil
		}
		return loadConfig()
	},
}

// loadConfig 解析工作区、加载配置文件并初始化 SQLite，注入 logic 层。
func loadConfig() error {
	if err := config.InitConf(workspace, cfgFile); err != nil {
		return err
	}
	logic.SetDB(config.Conf.DB)
	return nil
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径（优先级高于 -w）")
	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "工作区目录，默认 ~/.scan")
	// 关闭 cobra 默认注册的 completion 子命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
