package config

import (
	"fmt"
	"os"

	"github.com/Alexchent/goscan/cache"
	"github.com/spf13/viper"
)

type Config struct {
	Dir        string   // 保存日志文件的路径，不含文件
	FilterType []string // 需要过滤掉的文件类型，多个之间逗号分隔
	Cache      *Cache
	LogPath    string
	Sqlite     struct {
		DSN string
	}
}

type Cache struct {
	Addr     string
	Password string
	DB       int
}

var Conf = &Config{}
var FilterSuffix map[string]struct{}

// defaultConfig 返回默认配置
func defaultConfig() *Config {
	return &Config{
		Dir: "./scanLog",
		FilterType: []string{
			"js",
			"torrent",
			"downloading",
			"tiff",
		},
		Cache: &Cache{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}
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
	}

	// 验证必要配置
	if Conf.Cache == nil {
		return fmt.Errorf("配置缺少 Cache 配置项")
	}

	// 初始化过滤后缀映射
	FilterSuffix = make(map[string]struct{}, len(Conf.FilterType))
	for _, suffix := range Conf.FilterType {
		FilterSuffix[suffix] = struct{}{}
	}

	// 初始化 Redis 连接
	if err := cache.NewRedis(Conf.Cache.Addr, Conf.Cache.Password, Conf.Cache.DB); err != nil {
		return fmt.Errorf("redis 连接失败: %w", err)
	}

	return nil
}
