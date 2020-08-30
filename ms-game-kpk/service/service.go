package service 

import (
	"baihuatan/ms-game-kpk/model"
)

// import (
// 	"context"
// 	"log"
// )

// Service -
type Service interface {
	HealthCheck() bool
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

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service