package model

import (
	"encoding/json"
	"baihuatan/pkg/mysql"
	"log"
)

// ClientDetails 客户端描述
type ClientDetails struct {
	// client 标识
	ClientID      string
	// clinent 密钥
	ClientSecret   string
	// 访问令牌有效时间，秒
	AccessTokenValiditySeconds  int
	// 刷新令牌有效时间，秒
	RefreshTokenValiditySeconds  int
	// 重定向地址，授权码类型中使用
	RegisteredRedirectURI string
	// 可以使用的授权类型
	AuthorizedGrantTypes  []string
}

// IsMatch 是否匹配
func (clientDetails *ClientDetails) IsMatch(clientID string, clientSecret string) bool {
	return clientID == clientDetails.ClientID && clientSecret == clientDetails.ClientSecret
}

// ClientDetailsModel 对象
type ClientDetailsModel struct {

}

// NewClientDetailsModel 创建对象
func NewClientDetailsModel() *ClientDetailsModel {
	return &ClientDetailsModel{}
}

func (p *ClientDetailsModel) getTableName() string {
	return "client_details"
}

// GetClientDetailsByClientID 获取客户端信息
func (p *ClientDetailsModel) GetClientDetailsByClientID(clientID string) (*ClientDetails, error) {
	conn  := mysql.DB()
	if result, err := conn.Table(p.getTableName()).Where(map[string]interface{"client_id": clientID}).First(); err == nil {

		var authorizedGrantTypes [] string
		_ = json.Unmarshal([]byte(result["authorized_grant_types"].(string)), &authorizedGrantTypes)

		return &ClientDetails{
			ClientID:                   result["client_id"].(string),
			ClientSecret:               result["client_secret"].(string),
			AccessTokenValiditySeconds: int(result["access_token_validity_seconds"].(int64)),
			RefreshTokenValiditySeconds: int(result["refresh_token_validity_seconds"].(int64)),
			RegisteredRedirectURI: result["registered_redirect_uri"].(string),
			AuthorizedGrantTypes: authorizedGrantTypes,	
		}, nil
	}

	return nil, err
}

// CreateClientDetails 创建客户端详情
func (p *ClientDetailsModel) CreateClientDetails(clientDetails *ClientDetails) error {
	conn := mysql.DB()

	grantTypeString, _ := json.Marshal(clientDetails.AuthorizedGrantTypes)
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"client_id":      clientDetails.ClientID,
		"client_secret":  clientDetails.ClientSecret,
		"access_token_validity_seconds":  clientDetails.AccessTokenValiditySeconds,
		"refresh_token_validity_seconds": clientDetails.RefreshTokenValiditySeconds,
		"registered_redirect_uri":  clientDetails.RegisteredRedirectURI,
		"authorized_grant_types": grantTypeString,
	}).Insert()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	return nil
}
