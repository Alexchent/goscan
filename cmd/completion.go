package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 补全安装相关的标记。用块标记而非裸 source 行，便于精确判断是否已安装、
// 也方便未来批量管理（卸载/更新时按标记块操作）。
const (
	scanCompletionMarkBegin = "# >>> scan completion >>>"
	scanCompletionMarkEnd   = "# <<< scan completion <<<"
)

// installCompletion 在工作区目录 w 下生成 zsh 补全脚本文件 _scan（总是覆盖，保证内容
// 与当前子命令一致），并将其 source 行以标记块形式追加到用户 ~/.zshrc 中（若尚未安装）。
//
// 幂等：多次调用安全。
//   - _scan 文件每次覆盖刷新；
//   - .zshrc 仅在缺失标记块时追加一次，绝不重复。
//
// 失败不 panic，仅返回 error 由调用方决定是否中断流程。
func installCompletion(w string) error {
	completionFile := filepath.Join(w, "_scan")

	// 1. 生成补全脚本文件（覆盖）。
	// rootCmd 的所有子命令在包 init() 阶段已注册完毕，此处可直接生成。
	fd, err := os.OpenFile(completionFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建补全文件失败: %w", err)
	}
	defer fd.Close()
	if err := rootCmd.GenZshCompletion(fd); err != nil {
		return fmt.Errorf("生成补全脚本失败: %w", err)
	}
	fmt.Println("已生成补全脚本:", completionFile)

	// 2. 追加 source 行到 ~/.zshrc（若尚未安装）。
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}
	zshrc := filepath.Join(home, ".zshrc")

	if err := ensureZshrcBlock(zshrc, completionFile); err != nil {
		return fmt.Errorf("更新 .zshrc 失败: %w", err)
	}
	return nil
}

// ensureZshrcBlock 确保给定的 zshrc 中存在引用 completionFile 的标记块。
// 已存在则跳过；不存在则追加。
func ensureZshrcBlock(zshrc, completionFile string) error {
	if hasCompletionBlock(zshrc) {
		fmt.Println(".zshrc 已包含 scan 补全配置，跳过追加。")
		return nil
	}

	block := strings.Join([]string{
		"",
		scanCompletionMarkBegin,
		fmt.Sprintf(`[ -s "%s" ] && source "%s"`, completionFile, completionFile),
		scanCompletionMarkEnd,
		"",
	}, "\n")

	// 以追加模式打开；文件不存在时 O_CREATE 会创建。
	f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(block); err != nil {
		return err
	}
	fmt.Println("已添加 scan 补全配置到 .zshrc。")
	return nil
}

// hasCompletionBlock 检查 zshrc 中是否已存在 scan 补全标记块的起始标记。
func hasCompletionBlock(zshrc string) bool {
	f, err := os.Open(zshrc)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == scanCompletionMarkBegin {
			return true
		}
	}
	return false
}
