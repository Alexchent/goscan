package cmd

import (
	"fmt"
	"os"

	scan "github.com/Alexchent/goscan/service"
	"github.com/spf13/cobra"
)

// initCmd 创建/初始化工作区，生成默认配置文件。
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化工作区",
	Long: `初始化工作区：创建工作区目录并生成默认配置文件 .scan.yaml。

后续所有命令（start/find/clean/dump/load）都会自动使用工作区中的配置。
`,
	Run: func(cmd *cobra.Command, args []string) {
		w, err := MakeWorkSpace()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("工作区目录：", w)
		fmt.Println("工作区配置文件：", w+"/.scan.yaml")

		// 生成 zsh 补全脚本并追加到 ~/.zshrc（幂等）。补全安装失败不影响工作区初始化。
		// rootCmd.GenZshCompletion 作为生成函数传入，使 service 包无需反向依赖 cmd。
		if err := scan.InstallCompletion(w, rootCmd.GenZshCompletion); err != nil {
			fmt.Printf("⚠️ 补全安装失败（不影响工作区初始化）: %v\n", err)
		}
		fmt.Println("初始化完成")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func MakeWorkSpace() (string, error) {
	//home, _ := os.UserHomeDir()
	//workspace = filepath.Join(home, ".scan")
	// 判断目录是否存在
	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		// 如果目录不存在则创建目录
		if err := os.MkdirAll(workspace, 0755); err != nil {
			return "", fmt.Errorf("创建目录失败: %w", err)
		}
	}

	if err := scan.MakeConfig(workspace); err != nil {
		return "", err
	}

	return workspace, nil
}
