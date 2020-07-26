package config

import (
	"os"
	conf "baihuatan/pkg/config"
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
)

const (
	// KConfigType 配置类型
	KConfigType = "CONFIG_TYPE"
)

// Logger 日志
var Logger log.Logger

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)

	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}
	if err := conf.Sub("auth", &AuthPermitConfig); err != nil {
		Logger.Log("Fail to parse config", err)
	}
}

func initDefault()  {
	viper.SetDefault(KConfigType, "yaml")
}