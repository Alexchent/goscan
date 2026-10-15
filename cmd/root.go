package cmd

import (
	"fmt"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/service"
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
	// showVersion 由 -v/--version 触发，顶层 `scan -v` 时打印版本信息。
	showVersion bool
	keyword     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan 用于扫描文件目录，和查找文件",
	Long:  `scan 是一个用于扫描本地文件，形成日志文件，以便或许快速查找文件的程序`,
	// Run 仅在顶层无子命令时执行（如 `scan -v`）。
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			printVersion()
			return
		}
		if len(keyword) != 0 {
			count := 0
			count += service.Search(keyword)
			res := fmt.Sprintf("本次扫描发现 %d 个文件", count)
			fmt.Println(res)
			return
		}
		if scanDir != "" {
			startWork()
			return
		}
		// 无子命令且无 -v 时给出提示，等价于 -h 的简要说明。
		_ = cmd.Help()
	},
	// PersistentPreRunE 对所有子命令生效：先解析工作区并加载配置。
	// init 命令负责创建工作区，不要求工作区已初始化，故豁免。
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 顶层 `scan -v`：已初始化时展示完整配置，未初始化时降级为仅打印版本号。
		if showVersion {
			_ = loadConfig()
			return nil
		}
		// 以下命令不依赖工作区配置，豁免初始化校验：
		//   init      —— 负责创建工作区
		//   duplicate —— 纯内存扫描，不读写 DB
		//   remove    —— 读 JSON 执行 mv，不依赖工作区
		switch cmd.Name() {
		case "init", "duplicate", "remove":
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
	service.SetDB(config.Conf.DB)
	return nil
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径（优先级高于 -w）")
	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "工作区目录，默认 ~/.scan")
	// 顶层 `scan -v` 打印版本号，等价于 `scan version`
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "打印版本信息")
	rootCmd.Flags().StringVarP(&scanDir, "dir", "d", "", "需要扫描的目录")
	rootCmd.Flags().StringVarP(&keyword, "keywords", "k", "", "关键字")

	// 关闭 cobra 默认注册的 completion 子命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
