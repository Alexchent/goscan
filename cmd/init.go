package cmd

import (
	"fmt"

	"github.com/Alexchent/goscan/config"
	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
)

// initCmd 创建/初始化工作区，生成默认配置文件。
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化工作区",
	Long: `初始化工作区：创建工作区目录并生成默认配置文件 .scan.yaml。
默认工作区为 ~/.scan，可用 -w 指定其它目录。

后续所有命令（start/find/clean/dump/load）都会自动使用工作区中的配置。
可用 -w 切换到其它工作区，或用 -c 直接指定配置文件。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 确定目标工作区目录：优先 -w，否则默认 ~/.scan
		w := workspace
		if w == "" {
			home, err := config.UserHomeDir()
			if err != nil {
				fmt.Printf("获取用户目录失败: %v\n", err)
				return
			}
			w = home + "/.scan"
		}

		fmt.Println("工作区目录：", w)
		if err := scan.MakeConfig(w); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}
		fmt.Println("初始化完成。")
		fmt.Println("工作区配置文件：", w+"/.scan.yaml")
		fmt.Println("后续命令将自动使用该工作区，可用 -w 切换其它工作区。")

		// 生成 zsh 补全脚本并追加到 ~/.zshrc（幂等）。补全安装失败不影响工作区初始化。
		// rootCmd.GenZshCompletion 作为生成函数传入，使 service 包无需反向依赖 cmd。
		if err := scan.InstallCompletion(w, rootCmd.GenZshCompletion); err != nil {
			fmt.Printf("⚠️ 补全安装失败（不影响工作区初始化）: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
