package service

import (
	"baihuatan/ms-game-kpk/model"
	"context"
)

// import (
// 	"context"
// 	"log"
// )

// Service -
type Service interface {
	HealthCheck() bool
	GetUserData(context.Context, int64) (*model.KpkUser, error)
	AutoFetchQuestionsToCache(int64)
}

// KpkService -
type KpkService struct {
}

// HealthCheck -
func (o KpkService) HealthCheck() bool {
	return true
}

// AutoFetchQuestionsToCache - 
func (o KpkService) AutoFetchQuestionsToCache(num int64) {
	kpkModel := model.NewKpkQuestionModel()
	kpkModel.AutoFetchQuestionsToCache(num)
}



// GetUserData - 获取用户信息
func (o KpkService) GetUserData(ctx context.Context, userID int64) (*model.KpkUser, error) {
	return model.GetKpkUserByUID(ctx, userID)
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service