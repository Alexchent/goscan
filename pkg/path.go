package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetScanPath 获取扫描路径。
// scanDir 为命令行传入的目录参数，为空时交互式读取。
func GetScanPath(scanDir string) (string, error) {
	p := strings.TrimSpace(scanDir)
	if p == "" {
		fmt.Printf("请输入要扫描的目录: ")
		_, err := fmt.Scan(&p)
		if err != nil {
			return "", fmt.Errorf("读取输入失败: %w", err)
		}
		p = strings.TrimSpace(p)
	}

	// 处理特殊路径
	switch p {
	case "":
		//默认使用 Downloads 目录
		fmt.Println("警告: 不建议扫描根目录，将使用默认目录")
		dir, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("获取用户目录失败: %w", err)
		}
		return dir + "/Downloads", nil
	case ".":
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("获取当前目录失败: %w", err)
		}
		return wd, nil
	case "/":
		// 根目录扫描风险较大，
		return "", fmt.Errorf("根目录扫描风险较大，请重新选择")
	default:
		// 处理 ~ 开头的路径
		if strings.HasPrefix(p, "~") {
			dir, err := homedir.Dir()
			if err != nil {
				return "", fmt.Errorf("获取用户目录失败: %w", err)
			}
			p = strings.Replace(p, "~", dir, 1)
		}
		return strings.TrimRight(p, "/"), nil
	}
}

// ValidatePath 验证路径是否有效
func ValidatePath(p string) error {
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", p)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("无权限访问: %s", p)
		}
		return fmt.Errorf("无法访问路径: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("指定的路径不是目录: %s", p)
	}
	return nil
}
