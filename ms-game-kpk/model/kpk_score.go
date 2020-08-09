package model

import (
	"github.com/gohouse/gorose/v2"
	"baihuatan/pkg/mysql"
	"log"
)

// KpkScore 用户PK游戏积分
type KpkScore struct {
	UserID         int64  `json:"user_id"`   // 用户ID
	Score          int64  `json:"score"`     // 积分值
	UpdateTs       string  `json:"update_ts"` // 更新时间
}

// TableName -
func (t KpkScore) TableName() string {
	return "kpk_score"
}

// KpkScoreModel -
type KpkScoreModel struct {
}

// NewKpkScoreModel -
func NewKpkScoreModel() *KpkScoreModel {
	return &KpkScoreModel{}
}

// IncScore 增加积分
func (p *KpkScoreModel) IncScore(userID int64, score int64) bool {
	p.GetKpkScore(userID) //自动创建
	conn := mysql.DB()
	_, err := conn.Table(&KpkScore{}).Where("user_id", userID).Data(map[string]interface{}{"score": score}).Increment()
	
	if err == nil {
		return true
	}
	return false
}

// GetKpkScore
func (p *KpkScoreModel) GetKpkScore(userID int64) (*KpkScore, error) {
	conn := mysql.DB()
	sql := conn.Table(&KpkScore{})
	data, err := sql.Where("user_id", userID).First()
	if err != nil {
		log.Printf("Error: %v", err)
		// 记录不存在则创建
		_, err = sql.Data(map[string]interface{}{"user_id": userID, "score": 0}).Insert()
		if err != nil {
			log.Printf("Insert table kpk_score fail")
			return nil, err
		}
		// 再次获取
		data, err = sql.Where("user_id", userID).First()
		if err != nil {
			return nil, err
		}
	}

	return &KpkScore{
		UserID: data["user_id"].(int64),
		Score: data["score"].(int64),
		UpdateTs: data["update_ts"].(string),
	}, nil
}