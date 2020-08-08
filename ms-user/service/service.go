package service

import (
	"context"
	"baihuatan/ms-user/model"
	"log"
)

// Service Default a service interface
type Service interface {
	Check(ctx context.Context, username, password string) (int64, error)
	HealthCheck() bool
}

// UserService implement Service interface
type UserService struct {
}

// Check -
func (s UserService) Check(ctx context.Context, username string, password string) (int64, error) {
	userEntity := model.NewUserModel()
	res, err := userEntity.CheckUser(username, password)
	if err != nil {
		log.Printf("UserEntity.CreateUser, err : %v", err)
		return 0, err
	}
	return res.UserID, nil
}

// HealthCheck -
func (s UserService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service