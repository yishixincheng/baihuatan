package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// AdminEndpoints 端点
type PubEndpoints struct {
	HealthCheckEndpoint endpoint.Endpoint
}

// HealthCheck -- 健康检查
func (p *PubEndpoints) HealthCheck() bool {
	return true
}

// HealthRequest 健康检查请求结构
type HealthRequest struct {
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 健康检查响应结构
func MakeHealthCheckEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{Status: true}, nil
	}
}
