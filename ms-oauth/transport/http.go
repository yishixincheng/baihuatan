package transport

import (
	endpts "baihuatan/ms-oauth/endpoint"
	"baihuatan/ms-oauth/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// ErrBadRequest -
	ErrBadRequest = errors.New("invalid request parameter")
	// ErrGrantTypeRequest -
	ErrGrantTypeRequest = errors.New("invalid request grant type")
	// ErrTokenRequest -
	ErrTokenRequest = errors.New("invalid request token")
	// ErrInvalidClientRequest -
	ErrInvalidClientRequest = errors.New("无效客户端")
)

// MakeHTTPHandler make http handler use mux
func MakeHTTPHandler(ctx context.Context, endpoints endpts.OAuth2Endpoints,
	tokenService service.TokenService, clientService service.ClientDetailsService,
	zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, kitzipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}
	r.Path("metrics").Handler(promhttp.Handler())

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(clientService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("POST").Path("/token").Handler(kithttp.NewServer(
		endpoints.TokenEndpoint,
		decodeTokenRequest,
		encodeTokenResponse,
		clientAuthorizationOptions...,
	))

	r.Methods("POST").Path("/check_token").Handler(kithttp.NewServer(
		endpoints.CheckTokenEndpoint,
		decodeCheckTokenRequest,
		encodeJSONResponse,
		clientAuthorizationOptions...,
	))

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJSONResponse,
		options...,
	))

	return r
}

// encodeError

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// makeClientAuthorizationContext -
func makeClientAuthorizationContext(clientDetailsService service.ClientDetailsService, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		if clientID, clientSecret, ok := r.BasicAuth(); ok {
			fmt.Println("basic:" + clientID + ":" + clientSecret)
			clientDetails, err := clientDetailsService.GetClientDetailByClientID(ctx, clientID, clientSecret)
			fmt.Println(clientDetails)
			if err == nil {
				return context.WithValue(ctx, endpts.OAuth2ClientDetailsKey, clientDetails)
			}
		}
		return context.WithValue(ctx, endpts.OAuth2ErrorKey, ErrInvalidClientRequest)
	}
}

// decodeTokenRequest -
func decodeTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	grantType := r.URL.Query().Get("grant_type")
	if grantType == "" {
		return nil, ErrGrantTypeRequest
	}
	return &endpts.TokenRequest{
		GrantType: grantType,
		Reader:    r,
	}, nil
}

// decodeCheckTokenRequest -
func decodeCheckTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	tokenValue := r.URL.Query().Get("token")
	if tokenValue == "" {
		return nil, ErrTokenRequest
	}

	return &endpts.CheckTokenRequest{
		Token: tokenValue,
	}, nil
}

// encodeJSONResponse -
func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeTokenResponse - 401错误
func encodeTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if _, ok := ctx.Value(endpts.OAuth2ErrorKey).(error); ok {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest -
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}
