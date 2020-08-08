package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"baihuatan/ms-user/service"
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

// Check -
func (mw loggingMiddleware) Check(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"username", a,
			"pwd", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Check(ctx, a, b)
	return
}


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