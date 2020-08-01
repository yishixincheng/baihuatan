package transport

import (
	"context"
	endpts "baihuatan/oauth-service/endpoint"
)

// EncodeGPRCCheckTokenRequest -
func EncodeGRPCCheckTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpts.C)
}