package conf

import (
	"baihuatan/pkg/bootstrap"
	"baihuatan/pkg/discover"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
)

const (
	// ConfigType 配置类型
	ConfigType = "CONFIG_TYPE"
)

// ZipkinTracer 链路追踪
var ZipkinTracer *zipkin.Tracer

// Logger 日志
var Logger log.Logger

func initDefault() {
	viper.SetDefault(ConfigType, "yaml")
}

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	zipkinURL :=  "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Log("zipkin url", zipkinURL)
	initTracer(zipkinURL)
}

func initTracer(zipkinURL string) {
	var (
		err     error
		useNoopTracer = zipkinURL == ""
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)

	zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.HTTPConfig.Port)
	ZipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
	}
}

// LoadRemoteConfig 加载远程配置
func LoadRemoteConfig() (err error) {
	serviceInstance, err := discover.DiscoveryService(bootstrap.ConfigServerConfig.ID)
	if err != nil {
		return
	}
	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		configServer, bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName,
		bootstrap.ConfigServerConfig.Profile,
		viper.Get(ConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	viper.SetConfigType(viper.GetString(ConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	Logger.Log("Load config from: ", confAddr)
	return
}

// Sub 匹配
func Sub(key string, value interface{}) error {
	Logger.Log("配置文件的前缀为：", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}