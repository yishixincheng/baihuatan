package transport

import (
	endpts "baihuatan/ms-game-kpk/endpoint"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"baihuatan/ms-game-kpk/ws"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// ErrorBadRequest - 
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHTTPHandler - 路由
func MakeHTTPHandler(ctx context.Context, endpoints endpts.KpkEndpoints, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, kitzipkin.Name("http-transprt"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Path("metrics").Handler(promhttp.Handler())

	// 启动websocket长连接
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(ctx, ws.RoomM, w, r)
	})

	// 获取用户数据
	r.Methods("POST").Path("/getuserdata").Handler(kithttp.NewServer(
		endpoints.UserDataEndpoint,
		decodeGetUserDataRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/open/getuserdata").Handler(kithttp.NewServer(
		endpoints.UserDataEndpoint,
		decodeGetUserDataRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJSONResponse,
		options...,
	))

	return r
}


// encodeError -
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

// decodeGetUserDataRequest -
func decodeGetUserDataRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	getUserRequest := endpts.GetUserDataRequest{}
	if err := json.NewDecoder(r.Body).Decode(&getUserRequest); err != nil {
		return nil, err
	}
	if getUserRequest.UserID == 0 {
		userID, err := endpts.GetUserIDFromTokenParse(r)
		if err != nil {
			return nil, err
		}
		getUserRequest.UserID = userID
	}

	return &getUserRequest, nil
}

// decodeHealthCheckRequest -
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// encodeJSONResponse - 
func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}