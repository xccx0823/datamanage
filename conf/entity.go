package conf

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Settings struct {
	Dsn         string `yaml:"dsn"`
	MaxOpenCons int    `yaml:"maxOpenCons"`
	MaxIdleCons int    `yaml:"maxIdleCons"`

	WatchServer struct {
		DB struct {
			ServerId uint32 `yaml:"serverId"`
			Host     string `yaml:"host"`
			Port     uint16 `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Charset  string `yaml:"charset"`
		} `yaml:"db"`
		Kafka struct {
			Address        []string `yaml:"address"`
			FlushFrequency int      `yaml:"flushFrequency"`
		} `yaml:"kafka"`
		WithMonitorSyncTime int `yaml:"withMonitorSyncTime"`
	} `yaml:"watchServer"`
}

func InitConfig(configPath string) *Settings {
	envConfigPath := os.Getenv("DATA_MANAGE_CONFIG")
	if configPath == "" {
		if envConfigPath == "" {
			panic("没有配置数据管理工具所需配置，请在启动命令执行时添加 -c 指定配置文件，或者使用环境变量 DATA_MANAGE_CONFIG 指定")
		}
		configPath = envConfigPath
	}
	if !(strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml")) {
		panic("指定配置文件必须为yaml配置文件，以yaml或yml结尾")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	setting := new(Settings)
	if err := yaml.Unmarshal(data, setting); err != nil {
		panic(err)
	}
	return setting
}
