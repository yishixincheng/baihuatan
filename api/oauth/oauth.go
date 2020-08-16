package oauth

import (
	"baihuatan/api"
	"baihuatan/api/oauth/pb"
	"baihuatan/pkg/discover"
	"baihuatan/pkg/loadbalance"
	"context"

	"github.com/opentracing/opentracing-go"
)

// OAuthClient 鉴权
type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

// OAuthClientImpl 实现类
type OAuthClientImpl struct {
	manager		 api.ClientManager
	serviceName  string
	loadBalance  loadbalance.LoadBalance
	tracer       opentracing.Tracer
}

// CheckToken 检测令牌
func (impl *OAuthClientImpl) CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := impl.manager.DecoratorInvoke(ctx, "/pb.OAuthService/CheckToken", "token_check", tracer, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// NewOAuthClient 创建
func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = api.DefaultLoadBalance
	}

	return &OAuthClientImpl{
		manager: &api.DefaultClientManager{
			ServiceName: serviceName,
			LoadBalance: lb,
			DiscoveryClient: discover.ConsulService,
			Logger: discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil
}