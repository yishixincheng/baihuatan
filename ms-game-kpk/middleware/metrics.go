package middleware

import (
	"baihuatan/ms-game-kpk/service"
	"context"
	"time"
	"github.com/go-kit/kit/metrics"
)

// metricsMiddleware 定义监控中间件，嵌入Service
// 新增监控指标： requestCount 和 requestLatency
type metricsMiddleware struct {
	service.Service
	requestCount    metrics.Counter
	requestLatency  metrics.Histogram
}

// MetricsMiddleware 中间件
func MetricsMiddleware(requestCount metrics.Counter, requestLatency  metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return metricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

// HealthCheck -
func (mw metricsMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}