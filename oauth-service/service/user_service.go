package service

import (
	"baihuatan/oauth-service/model"
	"baihuatan/pb"
	"baihuatan/pkg/client"
	"context"
	"errors"
)

var (
	// ErrInvalidAuthentication 无效鉴权
	ErrInvalidAuthentication  = errors.New("Invalid auth")
	// ErrInvalidUserInfo 无效用户
	ErrInvalidUserInfo = errors.New("invalid user info")
)

// UserDetailsService 定义接口
type UserDetailsService interface {
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

// RemoteUserService struct
type RemoteUserService struct {
	userClient client.UserClient
}

// GetUserDetailByUsername 实现方法
func (service *RemoteUserService) GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	response, err := service.userClient.CheckUser(ctx, nil, &pb.UserRequest{
		Username: username,
		Password: password,
	})

	if err == nil {
		if response.UserID != 0 {
			return &model.UserDetails{
				UserID: response.UserID,
				Username: username,
				Password: password,
			}, nil
		}
		return nil, ErrInvalidUserInfo
	}
	return nil, err
}

// NewRemoteUserDetailService 创建对象
func NewRemoteUserDetailService() *RemoteUserService {
	userClient, _ := client.NewUserClient("user", nil, nil)
	return &RemoteUserService{
		userClient: userClient,
	}
}