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
	filename := filepath.Join(workspaceDir, ".scan.yaml")
	// 检查文件是否存在
	if pkg.FileIsExist(filename) {
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
	return nil
}
