package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Alexchent/goscan/cache"
	"github.com/Alexchent/goscan/service"
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

// InitConf 初始化配置，返回错误以便调用方处理
func InitConf(conf string) error {
	if len(conf) == 0 {
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败: %w", err)
		}
		conf = dir + "/.scan.yaml"
	}

	// 检查配置文件是否存在，不存在则自动创建
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		fmt.Printf("配置文件不存在，正在自动创建: %s\n", conf)
		confDir := filepath.Dir(conf)
		service.MakeConfig(confDir)
	}

	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证必要配置
	if Conf.Cache == nil {
		return fmt.Errorf("配置文件缺少 Cache 配置项")
	}

	// 初始化过滤后缀映射
	FilterSuffix = make(map[string]struct{}, len(Conf.FilterType))
	for _, suffix := range Conf.FilterType {
		FilterSuffix[suffix] = struct{}{}
	}

	// 初始化 Redis 连接
	if err := cache.NewRedis(Conf.Cache.Addr, Conf.Cache.Password, Conf.Cache.DB); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	return nil
}
