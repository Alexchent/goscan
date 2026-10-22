package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Alexchent/goscan/model"
	"github.com/spf13/viper"
)

type Config struct {
	FilterType []string `mapstructure:"FilterType"`
	Sqlite     struct {
		DSN string `mapstructure:"DSN"`
		// LogLevel 控制 GORM SQL 日志级别，取值：silent/error/warn/info。
		// 缺省 silent，不打印任何 SQL 日志（含慢 SQL 警告）。
		LogLevel string `mapstructure:"LogLevel"`
	} `mapstructure:"Sqlite"`
	// DB 在 InitConf 中初始化，持有共享的 SQLite 访问入口
	DB *model.ScanFile
}

var Conf = &Config{}
var FilterSuffix map[string]struct{}

// UserHomeDir 返回用户 Home 目录，对外暴露避免多处重复实现。
func UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// defaultConfig 返回默认配置（Dir 为空，表示用工作区目录）
func defaultConfig() *Config {
	return &Config{
		FilterType: []string{
			"js",
			"torrent",
			"downloading",
			"tiff",
		},
	}
}

// InitConf 加载配置文件并初始化 SQLite。
// workspace 为工作区目录（已解析为绝对路径），cfgFile 为可选的显式配置文件。
// 两者都为空时使用默认工作区 ~/.scan。若配置文件不存在则返回错误（要求先 scan init）。
func InitConf(workspace, cfgFile string) error {
	Conf = defaultConfig()

	// 确定工作区
	ws := workspace
	// 确定配置文件路径
	path := cfgFile
	if path == "" {
		path = filepath.Join(ws, ".scan.yaml")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s\n请先执行 scan init 初始化工作区", path)
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// DSN 为空时按 Dir 推导
	if Conf.Sqlite.DSN == "" {
		Conf.Sqlite.DSN = filepath.Join(ws, "scan.db")
	}

	// 初始化过滤后缀映射
	FilterSuffix = make(map[string]struct{}, len(Conf.FilterType))
	for _, suffix := range Conf.FilterType {
		FilterSuffix[suffix] = struct{}{}
	}

	// 初始化 SQLite（传入配置中的日志级别，缺省 silent）
	Conf.DB = model.NewScanFile(ws, Conf.Sqlite.LogLevel)
	return nil
}
