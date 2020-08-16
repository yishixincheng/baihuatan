package transport

import (
	endpts "baihuatan/ms-oauth/endpoint"
	"baihuatan/ms-oauth/model"
	"baihuatan/api/oauth/pb"
	"context"
)

// EncodeGRPCCheckTokenRequest -
func EncodeGRPCCheckTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpts.CheckTokenRequest)
	return &pb.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

// DecodeGRPCCheckTokenRequest -
func DecodeGRPCCheckTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CheckTokenRequest)
	return &endpts.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

// EncodeGRPCCheckTokenResponse - 
func EncodeGRPCCheckTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpts.CheckTokenResponse)
	
	if resp.Error != "" {
		return &pb.CheckTokenResponse{
			IsValidToken: false,
			Err:          resp.Error,
		}, nil
	}

	return &pb.CheckTokenResponse{
		UserDetails: &pb.UserDetails{
			UserID: resp.OAuthDetails.User.UserID,
			Username: resp.OAuthDetails.User.Username,
			Authorities: resp.OAuthDetails.User.Authorities,
		},
		ClientDetails: &pb.ClientDetails{
			ClientID: resp.OAuthDetails.Client.ClientID,
			AccessTokenValiditySeconds: int32(resp.OAuthDetails.Client.AccessTokenValiditySeconds),
			RefreshTokenValiditySeconds: int32(resp.OAuthDetails.Client.RefreshTokenValiditySeconds),
			AuthorizedGrantTypes: resp.OAuthDetails.Client.AuthorizedGrantTypes,
		},
		IsValidToken: true,
		Err:          "",
	},nil
}

// DecodeGRPCCheckTokenResponse -
func DecodeGRPCCheckTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.CheckTokenResponse)
	if resp.Err != "" {
		return endpts.CheckTokenResponse{
			Error: resp.Err,
		}, nil
	}
	return endpts.CheckTokenResponse{
		OAuthDetails: &model.OAuth2Details{
			User: &model.UserDetails{
				UserID: resp.UserDetails.UserID,
				Username: resp.UserDetails.Username,
				Authorities: resp.UserDetails.Authorities,
			},
			Client: &model.ClientDetails{
				ClientID:                    resp.ClientDetails.ClientID,
				AccessTokenValiditySeconds:  int(resp.ClientDetails.AccessTokenValiditySeconds),
				RefreshTokenValiditySeconds: int(resp.ClientDetails.RefreshTokenValiditySeconds),
				AuthorizedGrantTypes:        resp.ClientDetails.AuthorizedGrantTypes,
			},
		},
	}, nil
}