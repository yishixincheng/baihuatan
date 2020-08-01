package endpoint

import (
	"baihuatan/oauth-service/model"
	"baihuatan/oauth-service/service"
	"context"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// OAuth2Endpoints -
type OAuth2Endpoints struct {
	TokenEndpoint           endpoint.Endpoint
	CheckTokenEndpoint      endpoint.Endpoint
	GRPCCheckTokenEndpoint  endpoint.Endpoint
	HealthCheckEndpoint     endpoint.Endpoint
}

// MakeClientAuthorizationMiddleware -
func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); ok {
				return next(ctx, request)
			}
			return nil, ErrInvalidClientRequest
		}
	}
}

// MakeOAuth2AuthorizationMiddleware -
func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok {
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); ok{
				return next(ctx, request)
			}
			return nil, ErrInvalidUserRequest
		}
	}
}

// MakeAuthorityAuthorizationMiddleware -
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); ok{
				for _, value := range details.User.Authorities {
					if value == authority {
						return next(ctx, request)
					}
				}
				return nil, ErrNotPermit
			}

			return nil, ErrInvalidClientRequest
		}
	}
}

const (
	// OAuth2DetailsKey -
	OAuth2DetailsKey       = "OAuth2Details"
	// OAuth2ClientDetailsKey -
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	// OAuth2ErrorKey -
	OAuth2ErrorKey         = "OAuth2Error"
)

var (
	// ErrInvalidClientRequest -
	ErrInvalidClientRequest = errors.New("invalid client message")
	// ErrInvalidUserRequest -
	ErrInvalidUserRequest = errors.New("invalid user message")
	// ErrNotPermit -
	ErrNotPermit  = errors.New("not permit")
)

// TokenRequest -
type TokenRequest struct {
	GrantType  string
	Reader *http.Request
}

// TokenResponse -
type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error string `json:"error"`
}

// MakeTokenEndpoint -
func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*TokenRequest)
		token, err := svc.Grant(ctx, req.GrantType, ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails), req.Reader)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}

		return TokenResponse{
			AccessToken: token,
			Error: errString,
		}, nil
	}
}

// CheckTokenRequest -
type CheckTokenRequest struct {
	Token string
	ClientDetails model.ClientDetails
}

// CheckTokenResponse -
type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Error string `json:"error"`
}

// MakeCheckTokenEndpoint -
func MakeCheckTokenEndpoint(svc service.TokenService) 

// SimpleRequest - 
type SimpleRequest struct {
}

// SimpleReponse -
type SimpleReponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

// AdminRequest -
type AdminRequest struct {
}

// AdminResponse -
type AdminResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// HealthRequest -
type HealthRequest struct {
}
// HealthResponse -
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status:status,
		}, nil
	}
}


