package config

import (
	"fmt"
	"github.com/Alexchent/goscan/cache/mredis"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Dir        string   // 保存日志文件的路径，不含文件
	FilterType []string // 需要过滤掉的文件类型，多个之间逗号分隔
	Cache      *Cache
	LogPath    string
}

type Cache struct {
	Addr     string
	Password string
	DB       int
}

var Conf = &Config{}
var FilterSuffix map[string]struct{}

func InitConf(conf string) {
	if len(conf) == 0 {
		dir, _ := os.UserHomeDir()
		conf = dir + "/.scan.yaml"
	}
	viper.SetConfigFile(conf)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	FilterSuffix = make(map[string]struct{}, 0)
	for _, suffix := range Conf.FilterType {
		FilterSuffix[suffix] = struct{}{}
	}

	// 注册redis
	if Conf.Cache == nil {
		panic("redis 配置异常")
	}
	mredis.NewRedis(Conf.Cache.Addr, Conf.Cache.Password, Conf.Cache.DB)
}
