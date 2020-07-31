package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
	. "baihuatan/oauth-service/model"
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
	Grant(ctx context.Context, grantType string, client *ClientDetails, reader *http.Request) (*OAuth2Token, error)
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
func (tokenGranter *ComposeTokenGranter) Grant(ctx context.Context, grantType string, client *ClientDetails, reader *http.Request) (*OAuth2Token, error) {
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]
	
	if dispatchGranter == nil {
		return nil, ErrNotSupportGrantType
	}

	return dispatchGranter.Grant(ctx, grantType, client, reader)
}

// UsernamePasswordTokenGranter -
type UsernamePasswordTokenGranter struct {
	supportGrantType  string
}