package ratelimiter

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"github.com/juju/ratelimit"
	"time"
)

var (
	// ErrLimitExceed -
	ErrLimitExceed = errors.New("rate limit exceed")
)

// NewTokenBucketLimiterWithBuildIn - 使用x/time/rate创建限流中间件
func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// DynamicLimiter -
func DynamicLimiter(interval int, burst int) endpoint.Middleware {
	bucket := rate.NewLimiter(rate.Every(time.Second * time.Duration(interval)), burst)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bucket.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// NewTokenBucketLimiterWithJuju 使用juju/reatelimit创建中间件
func NewTokenBucketLimiterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}