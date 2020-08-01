package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
	"baihuatan/oauth-service/model"
)

var (
	// ErrNotSupportGrantType -
	ErrNotSupportGrantType = errors.New("grant type is not supported")
	// ErrNotSupportOperation -
	ErrNotSupportOperation = errors.New("no support operation")
	// ErrInvalidUsernameAndPasswordRequest -
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username, password")
	// ErrInvalidTokenRequest -
	ErrInvalidTokenRequest = errors.New("invalid token")
	// ErrExpiredToken -
	ErrExpiredToken = errors.New("token is expired")
)

// TokenGranter 令牌生成器
type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error)
}

// ComposeTokenGranter -
type ComposeTokenGranter struct {
	TokenGrantDict     map[string]TokenGranter
}

// NewComposeTokenGranter -
func NewComposeTokenGranter(tokenGranter map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{
		TokenGrantDict: tokenGranter,
	}
}

// Grant 生成令牌
func (tokenGranter *ComposeTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]
	
	if dispatchGranter == nil {
		return nil, ErrNotSupportGrantType
	}

	return dispatchGranter.Grant(ctx, grantType, client, reader)
}

// UsernamePasswordTokenGranter -
type UsernamePasswordTokenGranter struct {
	supportGrantType  string
	userDetailsService UserDetailsService
	tokenService TokenService
}

// NewUsernamePasswordTokenGranter -
func NewUsernamePasswordTokenGranter(grantType string, userDetailsService UserDetailsService, tokenService TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGrantType: grantType,
		userDetailsService: userDetailsService,
		tokenService: tokenService,
	}
}

// Grant - 生成令牌
func (tokenGranter *UsernamePasswordTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType {
		return nil, ErrNotSupportGrantType
	}
	
}


// TokenService -
type TokenService interface {
	// 根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	// 根据用户信息和客户端信息生成令牌
	GreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	// 根据用户信息和客户端信息获取已生成的访问令牌
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据访问令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
}