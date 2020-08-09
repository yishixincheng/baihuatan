package config

import (
	"github.com/go-kit/kit/log"
	"baihuatan/pkg/bootstrap"
	conf "baihuatan/pkg/config"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"os"
)

const (
	// KConfigType 配置类型
	KConfigType = "CONFIG_TYPE"
)

// ZipkinTracer 链路追踪
var ZipkinTracer *zipkin.Tracer
// Logger -
var Logger log.Logger

func init()  {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	if err := conf.Sub("mysql", &conf.MysqlConfig); err != nil {
		Logger.Log("Fail to parse mysql", err)
	}

	if err := conf.Sub("trace", &conf.TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	zipkinURL := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	Logger.Log("zipkin url", zipkinURL)

	initTracer(zipkinURL)
}

func initDefault()  {
	viper.SetDefault(KConfigType, "yaml")
}

func initTracer(zipkinURL string)  {
	var (
		err             error
		useNoopTracer   = zipkinURL == ""
		reporter        = zipkinhttp.NewReporter(zipkinURL)
	)
	// defer reporter.Close()
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