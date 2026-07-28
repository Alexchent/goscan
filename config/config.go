package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alexchent/goscan/model"
	"github.com/spf13/viper"
)

type Config struct {
	// Dir 保存日志文件与 scan.db 的目录路径
	Dir        string   `mapstructure:"Dir"`
	FilterType []string `mapstructure:"FilterType"`
	Sqlite     struct {
		DSN string `mapstructure:"DSN"`
	} `mapstructure:"Sqlite"`
	LogPath string `mapstructure:"LogPath"`
	// DB 在 InitConf 中初始化，持有共享的 SQLite 访问入口
	DB *model.ScanFile
}

var Conf = &Config{}
var FilterSuffix map[string]struct{}

// defaultConfig 返回默认配置
func defaultConfig() *Config {
	c := &Config{
		Dir: "./scanLog",
		FilterType: []string{
			"js",
			"torrent",
			"downloading",
			"tiff",
		},
	}
	c.Sqlite.DSN = filepath.Join(c.Dir, "scan.db")
	return c
}

// InitConf 初始化配置，返回错误以便调用方处理
func InitConf(c string) error {
	Conf = defaultConfig()
	// 检查配置文件是否存在
	if _, err := os.Stat(c); os.IsNotExist(err) {
		// 设置默认配置
		fmt.Printf("配置文件不存在，使用默认配置: %s\n", c)
	} else {
		// 配置文件存在，读取并合并
		viper.SetConfigFile(c)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
		// 将读取的配置信息覆盖默认配置
		if err := viper.Unmarshal(Conf); err != nil {
			return fmt.Errorf("解析配置文件失败: %w", err)
		}
		// 若配置未显式给出 DSN，则按 Dir 推导，保证两者一致
		if Conf.Sqlite.DSN == "" {
			Conf.Sqlite.DSN = filepath.Join(strings.TrimRight(Conf.Dir, "/"), "scan.db")
		}
	}

	// 确保目录存在，避免 gorm.Open 创建 db 文件失败
	if Conf.Dir != "" {
		if err := os.MkdirAll(Conf.Dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 初始化过滤后缀映射
	FilterSuffix = make(map[string]struct{}, len(Conf.FilterType))
	for _, suffix := range Conf.FilterType {
		FilterSuffix[suffix] = struct{}{}
	}

	// 初始化 SQLite（替换原 Redis 连接）
	// 注意：NewScanFile 接收的是目录路径，db 文件名为 scan.db
	dbDir := strings.TrimRight(Conf.Dir, "/")
	Conf.DB = model.NewScanFile(dbDir)

	return nil
}

// SqLPath 返回当前 SQLite 数据库文件的完整路径，供外部展示
func (c *Config) SqLPath() string {
	if c.Sqlite.DSN != "" {
		return c.Sqlite.DSN
	}
	if c.Dir == "" {
		return "scan.db"
	}
	return filepath.Join(strings.TrimRight(c.Dir, "/"), "scan.db")
}
