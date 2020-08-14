package middleware

import (
	"github.com/go-kit/kit/log"
	"baihuatan/ms-game-kpk/service"
	"time"
)

type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware 中间件
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}

// HealthCheck 健康检查
func (mw loggingMiddleware) HealthCheck() (ret bool) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "HealthCheck",
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret = mw.Service.HealthCheck()
	return
}