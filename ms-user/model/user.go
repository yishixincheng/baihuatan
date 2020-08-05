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
}

// UserModel -
type UserModel struct {
}

// NewUserModel -
func NewUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) getTableName() string {
	return "bht_user"
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
		UserName: data["user_name"].(string),
		Password: data["password"].(string),
		Birthday: data["birthday"].(string),
		Sex: int(data["sex"].(int64)),
	}

	return user, nil
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