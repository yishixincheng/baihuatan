package main

import (
	localconf "baihuatan/ms-oauth/config"
	endpts "baihuatan/ms-oauth/endpoint"
	"baihuatan/pkg/ratelimiter"
	"baihuatan/ms-oauth/service"
	"baihuatan/ms-oauth/transport"
	"baihuatan/pkg/bootstrap"
	conf "baihuatan/pkg/config"
	register "baihuatan/pkg/discover"
	"baihuatan/api/oauth/pb"
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
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	fmt.Println(bootstrap.HTTPConfig.Host)
	var (
		servicePort = flag.String("service.port", bootstrap.HTTPConfig.Port, "service port")
		grpcAddr = flag.String("grpc", bootstrap.RPCConfig.Port, "gRPC listen address.")
	)

	flag.Parse()

	fmt.Println(servicePort)
	fmt.Println(grpcAddr)

	ctx := context.Background()
	errChan := make(chan error)

	// 每秒钟产生100令牌
	ratebucket := rate.NewLimiter(rate.Every(time.Second * 1), 100)

	var (
		tokenService service.TokenService
		tokenGranter service.TokenGranter
		tokenEnhancer service.TokenEnhancer
		tokenStore service.TokenStore
		userDetailsService service.UserDetailsService
		clientDetailsService service.ClientDetailsService
		srv service.Service
	)

	// add logging middleware
	tokenEnhancer = service.NewJwtTokenEnhancer("secret")
	tokenStore = service.NewJwtTokenStore(tokenEnhancer.(*service.JwtTokenEnhancer))
	tokenService = service.NewTokenService(tokenStore, tokenEnhancer)
	userDetailsService = service.NewRemoteUserDetailService()
	clientDetailsService = service.NewMysqlClientDetailsService()
	srv = service.NewCommentService()

	tokenGranter = service.NewComposeTokenGranter(map[string]service.TokenGranter{
		"password": service.NewUsernamePasswordTokenGranter("password", userDetailsService, tokenService),
		"refresh_token": service.NewRefreshGranter("refresh_token", tokenService),
	})

	tokenEndpoint := endpts.MakeTokenEndpoint(tokenGranter, clientDetailsService)
	tokenEndpoint = endpts.MakeClientAuthorizationMiddleware(localconf.Logger)(tokenEndpoint)
	tokenEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(tokenEndpoint)
	tokenEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "token-endpoint")(tokenEndpoint)

	checkTokenEndpoint := endpts.MakeCheckTokenEndpoint(tokenService)
	checkTokenEndpoint = endpts.MakeClientAuthorizationMiddleware(localconf.Logger)(checkTokenEndpoint)
	checkTokenEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(checkTokenEndpoint)
	checkTokenEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "check-endpoint")(tokenEndpoint)

	gRPCCheckTokenEndpoint := endpts.MakeCheckTokenEndpoint(tokenService)
	gRPCCheckTokenEndpoint = ratelimiter.NewTokenBucketLimiterWithBuildIn(ratebucket)(gRPCCheckTokenEndpoint)
	gRPCCheckTokenEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "grpc-check-endpoint")(gRPCCheckTokenEndpoint)

	// 创建健康检查的Endpoint
	healthEndpoint := endpts.MakeHealthCheckEndpoint(srv)
	healthEndpoint = kitzipkin.TraceEndpoint(localconf.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpoints := endpts.OAuth2Endpoints{
		TokenEndpoint:       tokenEndpoint,
		CheckTokenEndpoint:  checkTokenEndpoint,
		GRPCCheckTokenEndpoint: gRPCCheckTokenEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHTTPHandler(ctx, endpoints, tokenService, clientDetailsService, localconf.ZipkinTracer, localconf.Logger)

	// http server
	go func() {
		fmt.Println("Http server start at port:" + *servicePort)
		mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User, conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
		// 启动前执行注册
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
		pb.RegisterOAuthServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <- errChan
	register.Deregister()
	fmt.Println(error)
}