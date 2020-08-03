package config

import (
	"github.com/spf13/viper"
)

var (
	cfgFile string = "config"
	cfgType string = "toml"
	cfgPath string = "/apps/conf/go/mlive.zhuanbo.gdxfhl.com/"
)

func Init() {
	viper.SetConfigName(cfgFile)
	viper.SetConfigType(cfgType)
	// 搜索目录 找到就不找了 可以多次调用
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath("./")
	// 适配mq子目录下面守护进程本地开发环境 by lemy 2020年03月16日14:09:14
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	if err := viper.ReadInConfig(); err != nil {
		panic("config read err: " + err.Error())
	}
}
