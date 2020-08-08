package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"baihuatan/ms-user/service"
)

var (
	// ErrInvalidRequestType -
	ErrInvalidRequestType = errors.New("invalid username,password")
)

// UserEndpoints 端点
type UserEndpoints struct {
	UserEndpoint           endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
}

// UserRequest -
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse -
type UserResponse struct {
	Result bool `json:"result"`
	UserID int64 `json:"user_id"`
	Error  error `json:"error"`
}

// Check -
func (u *UserEndpoints) Check(ctx context.Context, username, password string) (int64, error) {
	resp, err := u.UserEndpoint(ctx, UserRequest{
		Username: username,
		Password: password,
	})
	response := resp.(UserResponse)
	if (err != nil) {
		return 0, err
	}
	return response.UserID, nil
}

// HealthCheck - 健康检查
func (u *UserEndpoints) HealthCheck() bool {
	return true
}

// MakeUserEndpoint -
func MakeUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)
		username := req.Username
		password := req.Password

		userID, calError := svc.Check(ctx, username, password)
		if calError != nil {
			return UserResponse{Result: false, Error: calError}, nil
		}
		return UserResponse{Result: true, UserID: userID, Error: calError}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct {
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}