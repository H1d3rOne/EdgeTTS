package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config 是配置结构体
type Config struct {
	*viper.Viper
}

// NewConfig 创建并返回一个新的配置实例
func NewConfig() *Config {
	c := &Config{
		Viper: viper.New(),
	}

	// 设置配置文件名称
	c.SetConfigName("config")

	// 设置配置文件所在目录
	//c.AddConfigPath(".")
	c.AddConfigPath("../config")

	// 设置配置文件类型
	c.SetConfigType("yaml") // 或 "json", "toml", "hcl", "env", "properties"

	// 读取配置文件
	if err := c.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return c
}

var C = NewConfig() // 全局配置实例
