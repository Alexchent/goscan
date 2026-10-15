package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/Alexchent/goscan/config"
)

// version 在编译时通过 -ldflags "-X github.com/Alexchent/goscan/cmd.version=xxx" 注入，
// 默认值 dev。发布时由 Makefile / 构建脚本以 git tag 填入。
var version = "dev"

// currentVersion 返回当前版本号，去掉前导 v。
// 优先读取模块构建信息（go install module@tag 场景会自动带上解析的版本号），
// 其次用 ldflags 注入的 version，最后回退到 dev。
func currentVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return trimLeadingV(info.Main.Version)
	}
	return trimLeadingV(version)
}

// trimLeadingV 去掉版本号开头的 v（v1.0.4 -> 1.0.4）。
func trimLeadingV(s string) string {
	if len(s) > 0 && s[0] == 'v' {
		return s[1:]
	}
	return s
}

// printVersion 打印版本号，以及工作区初始化后完整的配置信息。
// 由顶层 `scan -v`（rootCmd.Run）调用。
func printVersion() {
	fmt.Println("当前版本是：v" + currentVersion())
	// 未初始化工作区时仅打印版本号
	if config.Conf.Workspace == "" || config.Conf.DB == nil {
		fmt.Println("(工作区未初始化，执行 scan init 初始化)")
		return
	}
	fmt.Println("工作区：", config.Conf.Workspace)
	fmt.Println("导出文件保存位置：", config.Conf.Dir)
	fmt.Println("数据库：", config.Conf.SqLPath())
	fmt.Println("扫描时忽略掉的文件类型：", config.Conf.FilterType)
}
