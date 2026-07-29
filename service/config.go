package service

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Alexchent/goscan/pkg"
)

//go:embed example.yaml
var configData string

// MakeConfig 在 workspaceDir 下生成默认配置文件 .scan.yaml。
// 若配置文件已存在则提示并跳过，不覆盖。
func MakeConfig(workspaceDir string) error {
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("创建工作区目录失败: %w", err)
	}

	filename := filepath.Join(workspaceDir, ".scan.yaml")
	if pkg.FileIsExist(filename) {
		fmt.Println("配置文件已存在:", filename)
		return nil
	}

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}
	defer fd.Close()

	if _, err := fd.WriteString(configData); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	fmt.Println("已生成配置文件:", filename)
	return nil
}
