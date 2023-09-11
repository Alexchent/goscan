package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Dir        string // 保存日志文件的路径，不含文件
	FilterType string // 需要过滤掉的文件类型，多个之间逗号分隔
}

var Conf = &Config{}

func init() {
	viper.SetConfigFile("./config/scan.yaml") // 指定配置文件路径
	err := viper.ReadInConfig()               // 读取配置信息
	if err != nil {                           // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	//fmt.Println(Conf)
}

func (c *Config) Get() *Config {
	return c
}
