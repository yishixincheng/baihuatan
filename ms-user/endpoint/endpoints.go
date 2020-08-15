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
	UserGetEndpoint        endpoint.Endpoint
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

// UserGetRequest -
type UserGetRequest struct {
	UserID int64 `json:"user_id"` 
}

// UserGetResponse -
type UserGetResponse struct {
	Result      bool      `json:"result"`
	Error       error     `json:"error"`
	UserID      int64     `json:"user_id"`
	UserName    string    `json:"user_name"` // 用户名称
	Birthday    string    `json:"birthday"`  // 生日
	Sex         int       `json:"sex"`       // 性别
	Avatar      string    `json:"avatar"`    // 头像
	City        string    `json:"city"`      // 城市
	District    string    `json:"district"`  // 区域
	Introduction string   `json:"introduction"`  // 介绍
	RoleID      int       `json:"role_id"`   // 性别
}

// MakeUserGetEndpoint -
func MakeUserGetEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserGetRequest)
		userID := req.UserID

		user, calError := svc.Get(ctx, userID)
		if calError != nil {
			return UserResponse{Result: false, Error: calError}, nil
		}
		return UserGetResponse{
				 Result: true,
				 Error:  nil,
				 UserID: userID,
				 UserName: user.UserName,
				 Birthday: user.Birthday,
				 Sex: user.Sex,
				 Avatar: user.Avatar,
				 City:   user.City,
				 District: user.District,
				 Introduction: user.Introduction,
				 RoleID: user.RoleID,
		}, nil
	}
}