package transport

import (
	"context"
	"baihuatan/api/user/pb"
	"errors"
	endpts "baihuatan/ms-user/endpoint"
)

// EncodeGRPCUserRequest -
func EncodeGRPCUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpts.UserRequest)
	return &pb.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

// DecodeGRPCUserRequest -
func DecodeGRPCUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserRequest)
	return endpts.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

// EncodeGRPCUserResponse -
func EncodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpts.UserResponse)
	if resp.Error != nil {
		return &pb.UserResponse{
			Result: bool(resp.Result),
			Err:  resp.Error.Error(),
		}, nil
	}

	return &pb.UserResponse{
		Result: bool(resp.Result),
		UserID: resp.UserID,
		Err: "",
	}, nil
}

// DecodeGRPCUserResponse -
func DecodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserResponse)
	return endpts.UserResponse{
		Result: bool(resp.Result),
		UserID: resp.UserID,
		Error:  errors.New(resp.Err),
	}, nil
}


// DecodeGRPCUserGetRequest -
func DecodeGRPCUserGetRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserGetRequest)
	return endpts.UserGetRequest{
		UserID: req.UserID,
	}, nil
}

// EncodeGRPCUserGetResponse -
func EncodeGRPCUserGetResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpts.UserGetResponse)
	
	if resp.Error != nil {
		return &pb.UserGetResponse{
			Result: bool(resp.Result),
			Err:  "error",
		}, nil
	}

	return &pb.UserGetResponse{
		Result: bool(resp.Result),
		Err: "",
		UserID: resp.UserID,
		UserName: resp.UserName,
		Sex: int32(resp.Sex),
		Birthday: resp.Birthday,
		City: resp.City,
		District: resp.District,
		Introduction: resp.Introduction,
		Avatar: resp.Avatar,
		RoleID: int32(resp.RoleID),
	}, nil
}