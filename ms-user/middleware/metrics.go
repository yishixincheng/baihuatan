package middleware

import (
	"baihuatan/ms-user/service"
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

// metricsMiddleware 定义监控中间件，嵌入Service
// 新增监控指标项： requestCount 和 requestLatency
type metricsMiddleware struct {
	service.Service
	requestCount    metrics.Counter
	requestLatency  metrics.Histogram
}

// MetricsMiddleware 封装监控方法
func MetricsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return metricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

// Check -
func (mw metricsMiddleware) Check(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Check"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.Check(ctx, a, b)
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
