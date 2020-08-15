package transport

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"baihuatan/pb"
	endpts "baihuatan/ms-user/endpoint"
)

type grpcServer struct {
	check grpc.Handler
	get   grpc.Handler
}

func (s *grpcServer) Check(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := s.check.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func(s *grpcServer) Get(ctx context.Context, r *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	_, resp, err := s.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserGetResponse), nil
}

// NewGRPCServer -
func NewGRPCServer(ctx context.Context, endpoints endpts.UserEndpoints, serverTracer grpc.ServerOption) pb.UserServiceServer {
	return &grpcServer{
		 check: grpc.NewServer(
			 endpoints.UserEndpoint,
			 DecodeGRPCUserRequest,
			 EncodeGRPCUserResponse,
			 serverTracer,
		 ),
		 get: grpc.NewServer(
			 endpoints.UserGetEndpoint,
			 DecodeGRPCUserGetRequest,
			 EncodeGRPCUserGetResponse,
			 serverTracer,
		 ),
	}
}