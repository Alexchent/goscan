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
}

type Cache struct {
	Addr     string
	Password string
	DB       int
}

var Conf = &Config{}

func InitConf() {
	u, _ := os.UserHomeDir()
	// 同时配置多个配置文件路径
	viper.AddConfigPath(u)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.SetConfigName("scan")
	viper.SetConfigType("yaml")
	//viper.SetConfigFile(u + "/scan.yaml") // 指定配置文件路径
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	// 注册redis
	fmt.Println(Conf.Cache)
	if Conf.Cache == nil {
		panic("redis 配置异常")
	}
	mredis.NewRedis(Conf.Cache.Addr, Conf.Cache.Password, Conf.Cache.DB)
}
