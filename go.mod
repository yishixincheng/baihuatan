module baihuatan

go 1.14

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/etcd-io/etcd v3.3.22+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/gohouse/gorose/v2 v2.1.7
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/hashicorp/consul/api v1.5.0
	github.com/juju/ratelimit v1.0.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/prometheus/client_golang v1.3.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.26.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
