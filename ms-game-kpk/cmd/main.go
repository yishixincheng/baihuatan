package main

import (
	localconf "baihuatan/ms-game-kpk/config"
	endpts "baihuatan/ms-game-kpk/endpoint"
	"baihuatan/ms-game-kpk/service"
	"baihuatan/ms-game-kpk/setup"
	"baihuatan/ms-game-kpk/transport"
	"baihuatan/ms-game-kpk/ws"
	"baihuatan/pkg/bootstrap"
	register "baihuatan/pkg/discover"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"time"

	"baihuatan/ms-game-kpk/middleware"
	"baihuatan/pkg/ratelimiter"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)


func main() {

	var (
		servicePort = flag.String("service.port", bootstrap.HTTPConfig.Port, "service port")
	//	grpcAddr    = flag.String("grpc", bootstrap.RPCConfig.Port, "gPRC listen address")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// 链接mysql，redis
	ws.InitWs()
	setup.InitMysql()
	setup.InitRedis()

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "bht",
		Subsystem: "ms_game_kpk",
		Name:      "request_count",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "bht",
		Subsystem: "ms_game_kpk",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	ratebucket := rate.NewLimiter(rate.Every(time.Second * 1), 100)

	var svc service.Service
	svc = service.KpkService{}

	// add logging middleware
	svc  = middleware.LoggingMiddleware(localconf.Logger)(svc)
	svc  = middleware.MetricsMiddleware(requestCount, requestLatency)(svc)

	// 健康检查Endpoint
	healthEndpoint := endpts.MakeHealthCheckEndpoint(svc)
	//healthEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(healthEndpoint)
	healthEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "health-endpoint")(healthEndpoint)

	// 用户信息获取
	userDataEndpoint := endpts.MakeUserDataEndpoint(svc)
	userDataEndpoint  = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(userDataEndpoint)
	userDataEndpoint  = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "get-user-data")(userDataEndpoint)

	endpoints := endpts.KpkEndpoints{
		HealthCheckEndpoint: healthEndpoint,
		UserDataEndpoint: userDataEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHTTPHandler(ctx, endpoints, localconf.ZipkinTracer, localconf.Logger)

	// http server
	go func() {
		fmt.Println("HTTP Server start at port:" + *servicePort)
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":" + *servicePort, handler)
	}()

	// grpc server
	// go func() {
	// 	fmt.Println("grpc Server start at port:" + *grpcAddr)
	// 	listener, err := net.Listen("tcp", ":" + *grpcAddr)
	// 	if err != nil {
	// 		errChan <- err
	// 		return
	// 	}
	// 	serverTracer := kitzipkin.GRPCServerTrace(localconf.ZipkinTracer, kitzipkin.Name("grpc-transport"))
	// 	tr := localconf.ZipkinTracer
	// 	md := metadata.MD{}
	// 	parentSpan := tr.StartSpan("test")

	// 	b3.InjectGRPC(&md)(parentSpan.Context())

	// 	ctx := metadata.NewIncomingContext(context.Background(), md)
	// 	handler := transport.NewGRPCServer(ctx, endpoints, serverTracer)
	// 	gRPCServer := grpc.NewServer()
	// 	pb.RegisterUserServiceServer(gRPCServer, handler)
	// 	errChan <- gRPCServer.Serve(listener)
	// }()

	// 死循环自动刷题到缓存中
	go func() {
		for {
			svc.AutoFetchQuestionsToCache(1000)
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// 服务注销
	error := <-errChan
	register.Deregister()
	fmt.Println(error)
}
