package model

import (
	"github.com/gohouse/gorose/v2"
	"baihuatan/pkg/mysql"
	"log"
)

// User -
type User struct {
	UserID      int64     `json:"user_id"`  // ID
	UserName    string    `json:"user_name"` // 用户名称
	Password    string    `json:"password"`  // 密码
	Birthday    string    `json:"birthday"`  // 生日
	Sex         int       `json:"sex"`       // 性别
	Avatar      string    `json:"avatar"`    // 头像
	City        string    `json:"city"`      // 城市
	District    string    `json:"district"`  // 区域
	Introduction string   `json:"introduction"`  // 介绍
	RoleID      int       `json:"role_id"`   // 性别
}

// UserModel -
type UserModel struct {
}

// NewUserModel -
func NewUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) getTableName() string {
	return "user"
}

// GetUserList - 
func (p *UserModel) GetUserList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

// CheckUser -
func (p *UserModel) CheckUser(username string, password string) (*User, error) {
	conn := mysql.DB()
	data, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"user_name": username, "password": password}).First()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	user := &User{
		UserID: data["user_id"].(int64),
	}

	return user, nil
}

// GetUser -
func (p *UserModel) GetUser(userID int64) (*User, error) {
	conn := mysql.DB()
	data, err := conn.Table(p.getTableName()).Where("user_id", userID).First()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return &User{
		UserID: data["user_id"].(int64),
		UserName: data["user_name"].(string),
		Birthday: data["birthday"].(string),
		Sex: int(data["sex"].(int32)),
		Avatar: data["avatar"].(string),
		City: data["city"].(string),
		District: data["district"].(string),
		Introduction: data["introduction"].(string),
		RoleID: int(data["role_id"].(int32)),
	}, nil
}

// CreateUser -
func (p *UserModel) CreateUser(user *User) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"user_id":  user.UserID,
		"user_name": user.UserName,
		"password": user.Password,
		"birthday": user.Birthday,
		"sex": user.Sex,
	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}