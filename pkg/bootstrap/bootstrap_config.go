package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	// _ "baihuatan/pkg/common"
	"log"
	"os"
)

func init() {
	viper.AutomaticEnv()
	initBootstrapConfig()
	//读取yaml文件
	
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	if err := subParse("http", &HTTPConfig); err != nil {
		log.Fatal("Fail to parse Http config", err)
	}
	if err := subParse("discover", &DiscoverConfig); err != nil {
		log.Fatal("Fail to parse Discover config", err)
	}
	if err := subParse("config", &ConfigServerConfig); err != nil {
		log.Fatal("Fail to parse config server", err)
	}

	if err := subParse("rpc", &RPCConfig); err != nil {
		log.Fatal("Fail to parse rpc server", err)
	}
}

func initBootstrapConfig() {
	// 设置读取的配置文件
	viper.SetConfigName("bootstrap")
	// 添加读取的配置文件路径
	viper.AddConfigPath("./")

	// windows环境下为%GOPATH，linux环境下为$GOPATH
	//configPath := common.GetGOPATH() + "/src"
    //viper.AddConfigPath(configPath)
	//设置文件类型
	viper.SetConfigType("yaml")
}

func subParse(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	if (sub == nil) {
		log.Printf("配置文件属性%v不存在", key)
		os.Exit(0)
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}