package endpoint

import (
	"context"
	"baihuatan/ms-game-kpk/service"
	"github.com/go-kit/kit/endpoint"
)

// KpkEndpoints 端点
type KpkEndpoints struct {
	HealthCheckEndpoint     endpoint.Endpoint
}

// HealthCheck - 健康检查
func (p *KpkEndpoints) HealthCheck() bool {
	return true
}

// HealthRequest 健康检查请求结构
type HealthRequest struct {
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status   bool    `json:"status"`
}

// MakeHealthCheckEndpoint 健康检查响应结构
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}