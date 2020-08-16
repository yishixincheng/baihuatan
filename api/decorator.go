package api

// 客户端装饰器
import (
	"baihuatan/pkg/bootstrap"
	conf "baihuatan/pkg/config"
	"baihuatan/pkg/discover"
	"baihuatan/pkg/loadbalance"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
)

var (
	// ErrRPCService rpc服务错误
	ErrRPCService = errors.New("no rpc service")
)

// DefaultLoadBalance - 
var DefaultLoadBalance loadbalance.LoadBalance = &loadbalance.WeightRoundRobinLoadBalance{}

// ClientManager 客户端管理
type ClientManager interface {
	DecoratorInvoke(ctx context.Context, path string, hystrixName string, tracer opentracing.Tracer,
	     inputVal interface{}, outVal interface{}) (err error)
}

// DefaultClientManager 客户端管理
type DefaultClientManager struct {
	ServiceName       string
	Logger            *log.Logger
	DiscoveryClient   discover.DiscoveryClient
	LoadBalance       loadbalance.LoadBalance
	after             []InvokerAfterFunc
	before            []InvokerBeforeFunc
}

// InvokerAfterFunc 后钩子
type InvokerAfterFunc  func() (err error)

// InvokerBeforeFunc 前钩子
type InvokerBeforeFunc  func() (err error)

// DecoratorInvoke 实现
func (manager *DefaultClientManager) DecoratorInvoke(ctx context.Context, path string, hystrixName string,
	tracer opentracing.Tracer, inputVal interface{}, outVal interface{}) (err error) {
	
	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}
	
	if err = hystrix.Do(hystrixName, func() error {
		instances := manager.DiscoveryClient.DiscoverServices(manager.ServiceName, manager.Logger)
		if instance, err := manager.LoadBalance.SelectService(instances); err == nil {
			if instance.GrpcPort > 0 {
				fmt.Println(instance.GrpcPort)
				if conn, err := grpc.Dial(instance.Host + ":" + strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
				    grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer))), grpc.WithTimeout(1*time.Second)); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						fmt.Println(err,"1")
						return err
					}
					fmt.Println("grpc:success")
				} else {
					fmt.Println(err,"2")
					return err
				}
			} else {
				return ErrRPCService
			}
		} else {
			return err
		}
		return nil
	}, func (e error) error {
		return e
	}); err != nil {
		return err
	}

	for _, fn := range manager.after {
		if err = fn(); err != nil {
			return err
		}
	}
	return nil
}

func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	zipkinURL := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	zipkinRecorder := bootstrap.HTTPConfig.Host + ":" + bootstrap.HTTPConfig.Port

	reporter := zipkinhttp.NewReporter(zipkinURL)
	//defer reporter.Close()

	// create our local service endpoint
	ept, err := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, zipkinRecorder)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(ept))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	res := zipkinot.Wrap(nativeTracer)
	// optionally set as Global OpenTracing tracer instance
	// opentracing.SetGlobalTracer(tracer)

	return res
}