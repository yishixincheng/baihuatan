package transport

import (
	"baihuatan/api/oauth/pb"
	endpts "baihuatan/ms-oauth/endpoint"
	"context"
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
			 DecodeGRPCCheckTokenRequest,
			 EncodeGRPCCheckTokenResponse,
			 serverTracer,
		),
	}
}