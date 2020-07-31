package service

import (
	"context"
	"errors"
	"baihuatan/oauth-service/model"
)

var (
	// ErrClientMessage 错误客户端
	ErrClientMessage = errors.New("invalid client")
)

// ClientDetailsService 定义接口
type ClientDetailsService interface {
	GetClientDetailByClientID(ctx context.Context, clientID, clientSecret string) (*model.ClientDetails, error)
}

// MysqlClientDetailsService 实现类
type MysqlClientDetailsService struct {
}

// NewMysqlClientDetailsService 实例化
func NewMysqlClientDetailsService() ClientDetailsService {
	return &MysqlClientDetailsService{}
}

// GetClientDetailByClientID 获取客户端
func (service *MysqlClientDetailsService)GetClientDetailByClientID(ctx context.Context, clientID, clientSecret string) (*model.ClientDetails, error) {
	clientDetailsModel := model.NewClientDetailsModel()
	var err error
	if clientDetails, err := clientDetailsModel.GetClientDetailsByClientID(clientID); err == nil {
		if clientSecret == clientDetails.ClientSecret {
			return clientDetails, nil
		} 
		return nil, ErrClientMessage
	}
	return nil, err
}
