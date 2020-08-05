package model

import (
	"time"
)

type OAuth2Token struct {
	// 刷新令牌
	RefreshToken  *OAuth2Token
	// 令牌类型
	TokenType  string
	// 令牌
	TokenValue  string
	// 过期时间
	ExpiresTime *time.Time
}

// IsExpired 令牌是否过期
func (oauth2Token *OAuth2Token) IsExpired() bool {
	return oauth2Token.ExpiresTime != nil && 
	    oauth2Token.ExpiresTime.Before(time.Now())
}

// OAuth2Details 基于oauth2协议
type OAuth2Details struct {
	Client    *ClientDetails
	User      *UserDetails
}
