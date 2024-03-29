package service

import (
	"baihuatan/ms-oauth/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrNotSupportGrantType -
	ErrNotSupportGrantType = errors.New("grant type is not supported")
	// ErrNotSupportOperation -
	ErrNotSupportOperation = errors.New("no support operation")
	// ErrNotEmptyUsernameAndPasswordRequest -
	ErrNotEmptyUsernameAndPasswordRequest = errors.New("username or password must required")
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
	TokenGrantDict map[string]TokenGranter
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

// TokenService -
type TokenService interface {
	// 根据访问令牌获取对应的用户信息和客户端信息
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	// 根据用户信息和客户端信息生成令牌
	CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据刷新令牌获取访问令牌
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	// 根据用户信息和客户端信息获取已生成的访问令牌
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根据访问令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
}

// UsernamePasswordTokenGranter -
type UsernamePasswordTokenGranter struct {
	supportGrantType   string
	userDetailsService UserDetailsService
	tokenService       TokenService
}

// NewUsernamePasswordTokenGranter -
func NewUsernamePasswordTokenGranter(grantType string, userDetailsService UserDetailsService, tokenService TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGrantType:   grantType,
		userDetailsService: userDetailsService,
		tokenService:       tokenService,
	}
}

// Grant - 生成令牌
func (tokenGranter *UsernamePasswordTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType {
		return nil, ErrNotSupportGrantType
	}
	// 从请求体中获取用户名和密码
	body, err := ioutil.ReadAll(reader.Body)
	fmt.Println(body, "请求体")
	if err != nil {
		return nil, ErrNotEmptyUsernameAndPasswordRequest
	}
	var paramKv = map[string]interface{}{}
	if err := json.Unmarshal(body, &paramKv); err != nil {
		return nil, ErrNotEmptyUsernameAndPasswordRequest
	}
	username := paramKv["username"].(string)
	password := paramKv["password"].(string)

	if username == "" || password == "" {
		return nil, ErrNotEmptyUsernameAndPasswordRequest
	}
	// 验证用户名密码是否正确
	fmt.Println("username and password:" + username + ":" + password)
	userDetails, err := tokenGranter.userDetailsService.GetUserDetailByUsername(ctx, username, password)

	if err != nil {
		return nil, err
		//return nil, ErrInvalidUsernameAndPasswordRequest
	}

	// 根据用户信息和客户端信息生成访问令牌
	return tokenGranter.tokenService.CreateAccessToken(&model.OAuth2Details{
		Client: client,
		User:   userDetails,
	})
}

// RefreshTokenGranter -
type RefreshTokenGranter struct {
	supportGrantType string
	tokenService     TokenService
}

// NewRefreshGranter -
func NewRefreshGranter(grantType string, tokenService TokenService) TokenGranter {
	return &RefreshTokenGranter{
		supportGrantType: grantType,
		tokenService:     tokenService,
	}
}

// Grant -
func (tokenGranter *RefreshTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType {
		return nil, ErrNotSupportGrantType
	}
	// 从请求中获取刷新令牌
	refreshTokenValue := reader.URL.Query().Get("refresh_token")

	if refreshTokenValue == "" {
		return nil, ErrInvalidTokenRequest
	}

	return tokenGranter.tokenService.RefreshAccessToken(refreshTokenValue)
}

// DefaultTokenService -
type DefaultTokenService struct {
	tokenStore    TokenStore
	tokenEnhancer TokenEnhancer
}

// NewTokenService 创建服务
func NewTokenService(tokenStore TokenStore, tokenEnhancer TokenEnhancer) TokenService {
	return &DefaultTokenService{
		tokenStore:    tokenStore,
		tokenEnhancer: tokenEnhancer,
	}
}

// CreateAccessToken -
func (tokenService *DefaultTokenService) CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	existToken, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
	var refreshToken *model.OAuth2Token
	if err == nil {
		// 存在未失效访问令牌，直接返回
		if !existToken.IsExpired() {
			tokenService.tokenStore.StoreAccessToken(existToken, oauth2Details)
			return existToken, nil
		}
		// 访问令牌已失效，移除
		tokenService.tokenStore.RemoveAccessToken(existToken.TokenValue)
		if existToken.RefreshToken != nil {
			refreshToken = existToken.RefreshToken
			tokenService.tokenStore.RemoveRefreshToken(refreshToken.TokenType)
		}
	}

	if refreshToken == nil || refreshToken.IsExpired() {
		refreshToken, err = tokenService.createRefreshToken(oauth2Details)
		if err != nil {
			return nil, err
		}
	}

	// 生成新的访问令牌
	accessToken, err := tokenService.createAccessToken(refreshToken, oauth2Details)
	if err == nil {
		// 保存新生成令牌
		tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
		tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
	}
	return accessToken, err
}

// createAccessToken 创建访问令牌
func (tokenService *DefaultTokenService) createAccessToken(refreshToken *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validitySeconds := oauth2Details.Client.AccessTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	expiredTime := time.Now().Add(s)
	accessToken := &model.OAuth2Token{
		RefreshToken: refreshToken,
		ExpiresTime:  &expiredTime,
		TokenValue:   uuid.NewV4().String(),
	}

	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(accessToken, oauth2Details)
	}
	return accessToken, nil
}

//  createRefreshToken -
func (tokenService *DefaultTokenService) createRefreshToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	validitySeconds := oauth2Details.Client.RefreshTokenValiditySeconds
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	expiredTime := time.Now().Add(s)
	refreshToken := &model.OAuth2Token{
		ExpiresTime: &expiredTime,
		TokenValue:  uuid.NewV4().String(),
	}

	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(refreshToken, oauth2Details)
	}

	return refreshToken, nil
}

// RefreshAccessToken -
func (tokenService *DefaultTokenService) RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error) {
	refreshToken, err := tokenService.tokenStore.ReadRefreshToken(refreshTokenValue)

	if err == nil {
		if refreshToken.IsExpired() {
			return nil, ErrExpiredToken
		}
		oauth2Details, err := tokenService.tokenStore.ReadOAuth2DetailsForRefreshToken(refreshTokenValue)
		// 移除原有访问令牌
		if err == nil {
			oauth2Token, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
			// 移除原有的访问令牌
			if err == nil {
				tokenService.tokenStore.RemoveAccessToken(oauth2Token.TokenValue)
			}

			// 移除已使用的刷新令牌
			tokenService.tokenStore.RemoveRefreshToken(refreshTokenValue)
			refreshToken, err = tokenService.createRefreshToken(oauth2Details)
			if err == nil {
				accessToken, err := tokenService.createAccessToken(refreshToken, oauth2Details)
				if err == nil {
					tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
					tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
				}
				return accessToken, err
			}

		}
	}
	return nil, err
}

// GetAccessToken -
func (tokenService *DefaultTokenService) GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return tokenService.tokenStore.GetAccessToken(details)
}

// ReadAccessToken -
func (tokenService *DefaultTokenService) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	return tokenService.tokenStore.ReadAccessToken(tokenValue)
}

// GetOAuth2DetailsByAccessToken -
func (tokenService *DefaultTokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error) {
	accessToken, err := tokenService.tokenStore.ReadAccessToken(tokenValue)
	if err == nil {
		if accessToken.IsExpired() {
			return nil, ErrExpiredToken
		}
		return tokenService.tokenStore.ReadOAuth2Details(tokenValue)
	}
	return nil, err
}

// TokenStore -
type TokenStore interface {
	// 存储访问令牌
	StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	// 根据令牌值获取访问令牌结构体
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
	// 根据令牌值获取令牌对应的客户端和用户信息
	ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error)
	// 根据客户端信息和用户信息获取访问令牌
	GetAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 移除存储的访问令牌
	RemoveAccessToken(tokenValue string)
	// 存储刷新令牌
	StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	// 移除存储的刷新令牌
	RemoveRefreshToken(tokenType string)
	// 根据令牌值获取刷新令牌
	ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error)
	// 根据令牌值获取刷新令牌对应的客户端和用户信息
	ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error)
}

// NewJwtTokenStore -
func NewJwtTokenStore(jwtTokenEnhancer *JwtTokenEnhancer) TokenStore {
	return &JwtTokenStore{
		jwtTokenEnhancer: jwtTokenEnhancer,
	}
}

// JwtTokenStore -
type JwtTokenStore struct {
	jwtTokenEnhancer *JwtTokenEnhancer
}

// StoreAccessToken -
func (tokenStore *JwtTokenStore) StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
	tokenStore.jwtTokenEnhancer.Enhance(oauth2Token, oauth2Details)
}

// ReadAccessToken -
func (tokenStore *JwtTokenStore) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err
}

// ReadOAuth2Details 根据令牌获取对应的客户端信息
func (tokenStore *JwtTokenStore) ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error) {
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err
}

// GetAccessToken 根据客户端信息和用户信息获取访问令牌
func (tokenStore *JwtTokenStore) GetAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return nil, ErrNotSupportOperation
}

// RemoveAccessToken 移除存储的访问令牌
func (tokenStore *JwtTokenStore) RemoveAccessToken(tokenValue string) {

}

// StoreRefreshToken 存储刷新令牌
func (tokenStore *JwtTokenStore) StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {

}

// RemoveRefreshToken 移除存储的刷新令牌
func (tokenStore *JwtTokenStore) RemoveRefreshToken(tokenType string) {

}

// ReadRefreshToken 根据令牌获取刷新令牌
func (tokenStore *JwtTokenStore) ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error) {
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err
}

// ReadOAuth2DetailsForRefreshToken 根据令牌获取客户端信息
func (tokenStore *JwtTokenStore) ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error) {
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err
}

// TokenEnhancer -
type TokenEnhancer interface {
	// 组装 Token 信息
	Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 从 Token 中还原信息
	Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error)
}

// OAuth2TokenCustomClaims -
type OAuth2TokenCustomClaims struct {
	UserDetails   model.UserDetails
	ClientDetails model.ClientDetails
	RefreshToken  model.OAuth2Token
	jwt.StandardClaims
}

// JwtTokenEnhancer -
type JwtTokenEnhancer struct {
	secretKey []byte
}

// Enhance - 组装Token信息
func (enhancer *JwtTokenEnhancer) Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return enhancer.sign(oauth2Token, oauth2Details)
}

// Extract 解包
func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &OAuth2TokenCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey, nil
	})

	if err == nil {
		claims := token.Claims.(*OAuth2TokenCustomClaims)
		expiresTime := time.Unix(claims.ExpiresAt, 0)

		return &model.OAuth2Token{
				RefreshToken: &claims.RefreshToken,
				TokenValue:   tokenValue,
				ExpiresTime:  &expiresTime,
			}, &model.OAuth2Details{
				User:   &claims.UserDetails,
				Client: &claims.ClientDetails,
			}, nil
	}
	return nil, nil, err
}

// sign 签名
func (enhancer *JwtTokenEnhancer) sign(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	expireTime := oauth2Token.ExpiresTime
	clientDetails := *oauth2Details.Client
	userDetails := *oauth2Details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	claims := OAuth2TokenCustomClaims{
		UserDetails:   userDetails,
		ClientDetails: clientDetails,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "System",
		},
	}

	if oauth2Token.RefreshToken != nil {
		claims.RefreshToken = *oauth2Token.RefreshToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenValue, err := token.SignedString(enhancer.secretKey)

	if err == nil {
		oauth2Token.TokenValue = tokenValue
		oauth2Token.TokenType = "jwt"
		return oauth2Token, nil
	}
	return nil, err
}

// NewJwtTokenEnhancer -
func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: []byte(secretKey),
	}
}
