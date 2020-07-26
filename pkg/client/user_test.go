package client

import (
	"baihuatan/pb"
	"context"
	"fmt"
	"testing"
	"log"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func TestUserClientImpl_CheckUser(t *testing.T) {
	client, _ := NewUserClient("user", nil, genTracerAct(nil))

	if response, err := client.CheckUser(context.Background(), nil, &pb.UserRequest{
		Username: "testuser",
		Password: "123456",
	}); err == nil {
		fmt.Println(response.Result)
	} else {
		fmt.Println(err.Error())
	}
}

func genTracerAct(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	zipkinURL := "http://192.168.43.251:9411/api/v2/spans"
	zipkinRecorder := "localhost:12344"

	reporter := zipkinhttp.NewReporter(zipkinURL)
	defer reporter.Close()

	ept, err := zipkin.NewEndpoint("user-client", zipkinRecorder)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(ept))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	return zipkinot.Wrap(nativeTracer)

}