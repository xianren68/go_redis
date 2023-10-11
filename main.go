package main

import (
	"fmt"
	"go_redis/config"
	"go_redis/lib/logger"
	"go_redis/tcp"
	"os"
)

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6399,
}

// 判断文件是否存在
func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	// 没有错误并且不是文件夹
	return err == nil && !stat.IsDir()
}

// 默认配置

func main() {
	// 设置日志格式
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	// 读取配置文件
	configName := "redis.config"
	if fileExists(configName) {
		config.SetupConfig(configName)
	} else {
		// 使用默认配置
		config.Properties = defaultProperties
	}

	// 开启tcp连接
	echoHandler := tcp.MakeEchoHandler()
	err := tcp.ListenAndServerWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	},
		echoHandler)
	if err != nil {
		logger.Error(err)
	}
}
