package user

import (
	"baihuatan/api"
	"baihuatan/api/user/pb"
	"baihuatan/pkg/discover"
	"baihuatan/pkg/loadbalance"
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

// UserClient 接口
type UserClient interface {
	CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error)
	GetUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserGetRequest) (*pb.UserGetResponse, error)
}

// UserClientImpl 实现类
type UserClientImpl struct {
	/**
	* 可以配置负载均衡策略，重试、等机制。也可以配置invokeAfter和invokerBefore
	*/
	manager       api.ClientManager
	serviceName   string
	loadBlance    loadbalance.LoadBalance
	tracer        opentracing.Tracer
}

// CheckUser 验证用户
func (impl *UserClientImpl) CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	if err := impl.manager.DecoratorInvoke(ctx, "/pb.UserService/Check", "user_check", tracer, request, response); err != nil {
		return  nil, err
	}
	fmt.Println(response, "grpc接口返回")
	return response, nil	
}

// GetUser 获取用户信息
func (impl *UserClientImpl) GetUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	response := new(pb.UserGetResponse)
	if err := impl.manager.DecoratorInvoke(ctx, "/pb.UserService/Get", "user_get", tracer, request, response); err != nil {
		return nil, err
	}
	fmt.Println(response, "grpc接口返回")
	return response, nil
} 

// NewUserClient 创建实例
func NewUserClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (UserClient, error) {
	if serviceName == "" {
		serviceName = "user"
	}
	if lb == nil {
		lb = api.DefaultLoadBalance
	}

	return &UserClientImpl{
		manager: &api.DefaultClientManager{
			ServiceName: serviceName,
			LoadBalance: lb,
			DiscoveryClient: discover.ConsulService,
			Logger: discover.Logger,
		},
		serviceName: serviceName,
		loadBlance: lb,
		tracer:     tracer,
	}, nil
}