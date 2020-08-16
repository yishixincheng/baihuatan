package main

import (
	localconf "baihuatan/ms-user/config"
	endpts "baihuatan/ms-user/endpoint"
	"baihuatan/ms-user/service"
	"baihuatan/ms-user/transport"
	"baihuatan/api/user/pb"
	"baihuatan/pkg/bootstrap"
	conf "baihuatan/pkg/config"
	register "baihuatan/pkg/discover"
	"baihuatan/pkg/mysql"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"baihuatan/ms-user/middleware"
	"baihuatan/pkg/ratelimiter"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {

	var (
		servicePort = flag.String("service.port", bootstrap.HTTPConfig.Port, "service port")
		grpcAddr = flag.String("grpc", bootstrap.RPCConfig.Port, "gRPC listen address.")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "bht",
		Subsystem: "ms_user",
		Name:      "request_count",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "bht",
		Subsystem: "ms_user",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	ratebucket := rate.NewLimiter(rate.Every(time.Second * 1), 100)

	var svc service.Service
	svc = service.UserService{}

	// add logging middleware
	svc = middleware.LoggingMiddleware(localconf.Logger)(svc)
	svc = middleware.MetricsMiddleware(requestCount, requestLatency)(svc)

	userEndpoint := endpts.MakeUserEndpoint(svc)
	userEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(userEndpoint)
	userEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "user-endpoint")(userEndpoint)

	// 健康检查Endpoint
	healthEndpoint := endpts.MakeHealthCheckEndpoint(svc)
	healthEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpoints := endpts.UserEndpoints{
		UserEndpoint: userEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHTTPHandler(ctx, endpoints, localconf.ZipkinTracer, localconf.Logger)

	// http server
	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		mysql.InitMysql(conf.MysqlConfig.Host, 
						conf.MysqlConfig.Port,
						conf.MysqlConfig.User,
						conf.MysqlConfig.Pwd,
						conf.MysqlConfig.Db)
		// 注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":" + *servicePort, handler)				
	}()

	// grpc server
	go func() {
		fmt.Println("grpc Server start at port:" + *grpcAddr)
		listener, err := net.Listen("tcp", ":" + *grpcAddr)
		if err != nil {
			errChan <- err
			return
		}
		serverTracer := kitzipkin.GRPCServerTrace(localconf.ZipkinTracer, kitzipkin.Name("grpc-transport"))
		tr := localconf.ZipkinTracer
		md := metadata.MD{}
		parentSpan := tr.StartSpan("test")

		b3.InjectGRPC(&md)(parentSpan.Context())

		ctx := metadata.NewIncomingContext(context.Background(), md)
		handler := transport.NewGRPCServer(ctx, endpoints, serverTracer)
		gRPCServer := grpc.NewServer()
		pb.RegisterUserServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
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