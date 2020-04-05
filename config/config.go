package config

import (
	"github.com/spf13/viper"
	"strings"
)

var Values *Config

type User struct {
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

type Config struct {
	Name string

	Log struct {
		Writers       string `mapstructure:"writers"`
		LoggerLevel   string `mapstructure:"logger_level"`
		LoggerFile    string `mapstructure:"logger_file"`
		LogFormatText bool   `mapstructure:"log_format_text"`
	} `mapstructure:"log"`

	Mysql struct {
		Addr string `mapstructure:"addr"` //ip地址端口号
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
		DB   string `mapstructure:"db"`
	} `mapstructure:"mysql"`

	Users []User `mapstructure:"users"`

	Servers []string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	Values = &c

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("./") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")          // 设置配置文件格式为YAML
	viper.AutomaticEnv()                 // 读取匹配的环境变量
	viper.SetEnvPrefix("MATRIX_SERVICE") // 读取环境变量的前缀为VEHICLE_MONITORING
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}
