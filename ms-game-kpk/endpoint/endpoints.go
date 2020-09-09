package endpoint

import (
	"baihuatan/api/oauth/pb"
	"context"
	"baihuatan/ms-game-kpk/service"
	"github.com/go-kit/kit/endpoint"
	"baihuatan/ms-game-kpk/model"
	"net/http"
	"encoding/json"
)

// KpkEndpoints 端点
type KpkEndpoints struct {
	HealthCheckEndpoint     endpoint.Endpoint
	UserDataEndpoint        endpoint.Endpoint
}

// HealthCheck - 健康检查
func (p *KpkEndpoints) HealthCheck() bool {
	return true
}

// HealthRequest 健康检查请求结构
type HealthRequest struct {
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status   bool    `json:"status"`
}

// MakeHealthCheckEndpoint 健康检查响应结构
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}

// GetUserDataRequest -
type GetUserDataRequest struct {
	UserID     int64
}

// GetUserDataResponse -
type GetUserDataResponse struct {
	UserData  *model.KpkUser  `json:"user_data"`
	Error     string          `json:"error"`

}

// MakeUserDataEndpoint 用户数据获取
func MakeUserDataEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetUserDataRequest)

		kpkUser, err := svc.GetUserData(ctx, req.UserID)

		if err != nil {
			return GetUserDataResponse{
				nil, 
				err.Error(),
			}, err
		}

		return GetUserDataResponse {
			UserData: kpkUser,
			Error: "",
		}, nil
	}
}

// GetUserIDFromTokenParse -
func GetUserIDFromTokenParse(r *http.Request) (int64, error) {
	userToken := r.Header.Get("Bhtuser")
	var userDetail = pb.UserDetails{}
	if err := json.Unmarshal([]byte(userToken), &userDetail); err != nil {
		return 0, err
	}

	return userDetail.UserID, nil
}