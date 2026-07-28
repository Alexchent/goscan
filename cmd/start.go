package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexchent/goscan/config"
	"github.com/Alexchent/goscan/logic"
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "开始扫描",
	Long:  `开始扫描`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			fmt.Println("扫描完成，耗时：", time.Since(start))
			// 确保程序退出时关闭日志文件
			_ = logic.CloseLogfile()
		}()

		// 获取扫描路径
		scanPath, err := getScanPath()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		// 验证路径
		if err := validatePath(scanPath); err != nil {
			fmt.Printf("路径验证失败: %v\n", err)
			return
		}

		color.HiGreen.Println("开始扫描：", scanPath)
		if err := logic.Save(logic.CacheKey, scanPath, config.FilterSuffix); err != nil {
			fmt.Printf("扫描出错: %v\n", err)
		}
	},
}

// getScanPath 获取扫描路径
func getScanPath() (string, error) {
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

// validatePath 验证路径是否有效
func validatePath(p string) error {
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

func init() {
	startCmd.Flags().StringVarP(&scanDir, "dir", "d", "", "需要扫描的目录")
	rootCmd.AddCommand(startCmd)
}
