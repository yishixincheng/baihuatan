package service 

// import (
// 	"context"
// 	"log"
// )

// Service -
type Service interface {
	HealthCheck() bool
}

// KpkService -
type KpkService struct {
}

// HealthCheck -
func (o KpkService) HealthCheck() bool {
	return true
}