package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpts "baihuatan/ms-user/endpoint"
	"github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// ErrorBadRequest -
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHTTPHandler make http handler use mux
func MakeHTTPHandler(ctx context.Context, endpoints endpts.UserEndpoints, zipkinTracer *zipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, kitzipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("POST").Path("/check/valid").Handler(kithttp.NewServer(
		endpoints.UserEndpoint,
		decodeUserRequest,
		encodeUserResponse,
		options...,	
	))

	r.Path("metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeUserResponse,
		options...,
	))

	return r
}

// decodeUserRequest -
func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userRequest endpts.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return nil, err
	}
	return userRequest, nil
}

// encodeUserResponse -
func encodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest -
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
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