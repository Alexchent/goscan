package cmd

import (
	"fmt"

	"github.com/Alexchent/goscan/config"
	"github.com/spf13/cobra"
)

// version 在编译时通过 -ldflags "-X github.com/Alexchent/goscan/cmd.version=xxx" 注入，
// 默认值 dev。发布时由 Makefile / 构建脚本以 git tag 填入。
var version = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is goscan`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("当前版本是：v" + version)
		// 未初始化工作区时仅打印版本号
		if config.Conf.Workspace == "" || config.Conf.DB == nil {
			fmt.Println("(工作区未初始化，执行 scan init 初始化)")
			return
		}
		fmt.Println("工作区：", config.Conf.Workspace)
		fmt.Println("导出文件保存位置：", config.Conf.Dir)
		fmt.Println("数据库：", config.Conf.SqLPath())
		fmt.Println("扫描时忽略掉的文件类型：", config.Conf.FilterType)
	},
}
