package transport

import (
	"baihuatan/pb"
	endpts "baihuatan/oauth-service/endpoint"
	"context"
	"baihuatan/oauth-service/model"
	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	checkTokenServer grpc.Handler
}

func (s *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, resp, err := s.checkTokenServer.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CheckTokenResponse), nil
}

// NewGRPCServer -
func NewGRPCServer(ctx context.Context, endpoints endpts.OAuth2Endpoints, serverTracer grpc.ServerOption) pb.OAuthServiceServer{
	return &grpcServer{
		checkTokenServer: grpc.NewServer(
			 endpoints.GRPCCheckTokenEndpoint,
			 DecodeGRPC

		),
	}
}

