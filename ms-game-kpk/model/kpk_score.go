package model

import (
	"baihuatan/api/user"
	"baihuatan/api/user/pb"
	conf "baihuatan/pkg/config"
	"baihuatan/pkg/mysql"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/go-redis/redis"
)

// KpkScore 用户PK游戏积分
type KpkScore struct {
	UserID         int64      `gorose:"user_id" json:"user_id"`   // 用户ID
	Score          int64      `gorose:"score" json:"score"`     // 积分值
	PetID          int64      `gorose:"pet_id" json:"pet_id"`    // 宠物ID
	RoadID         int64      `gorose:"road_id" json:"road_id"`   // 跑道ID
	UpdateTs       time.Time  `gorose:"update_ts" json:"update_ts"` // 更新时间
	Rank           string     `json:"rank"`                         //等级
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
		// 清除缓存
		ClearKpkUserCache(userID)
		return true
	}
	return false
}

// BatchIncScore 批量增加积分
func (p *KpkScoreModel) BatchIncScore(data []map[string]interface{}) {
	for _, item := range data {
		p.IncScore(item["userID"].(int64), item["score"].(int64))
	}
}

// GetKpkScore -
func (p *KpkScoreModel) GetKpkScore(userID int64) (*KpkScore, error) {
	conn := mysql.DB()
	sql := conn.Table(&KpkScore{})
	data, err := sql.Where("user_id", userID).First()
	if err != nil || data == nil {
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
		PetID: data["pet_id"].(int64),
		RoadID: data["road_id"].(int64),
		UpdateTs: data["update_ts"].(time.Time),
		Rank: GetRankByScore(data["score"].(int64)),
	}, nil
}

// KpkUser -用户信息
type KpkUser struct{
	KpkScore
	UserName    string    `gorose:"user_name" json:"user_name"`
	Sex         int32     `gorose:"sex" json:"sex"`
	Avatar      string    `gorose:"avatar" json:"avatar"`
	Birthday    string    `gorose:"birthday" json:"birthday"`
	City        string    `gorose:"city" json:"city"`
	District    string    `gorose:"district" json:"district"`
	RoleID      int32     `gorose:"role_id" json:"role_id"`
}


// GetKpkUserByUID - 
func GetKpkUserByUID(ctx context.Context, UID int64) (*KpkUser, error){
	conn := conf.Redis.RedisConn
	userKey := fmt.Sprintf("KpkUser_%v", UID)
	kpkUserStr, err := conn.Get(userKey).Result()
    if err == redis.Nil {
		// key 不存在，创建
		apiClient, _ := user.NewUserClient("user", nil, nil)
		if resp, apiErr := apiClient.GetUser(ctx, nil, &pb.UserGetRequest{
			UserID: UID,
		}); apiErr != nil {
			return nil, apiErr
		} else {
			// 找到用户
			kpkScoreModel := NewKpkScoreModel()
			kpkScore, err2 := kpkScoreModel.GetKpkScore(UID);
			if err2 != nil {
				return nil, err2
			}

			kpkUser := &KpkUser{
				KpkScore: *kpkScore,
				UserName: resp.UserName,
				Sex:      resp.Sex,
				Avatar:   resp.Avatar,
				City:     resp.City,
				Birthday: resp.Birthday,
				District: resp.District,
				RoleID:   resp.RoleID,
			}

			// 保存到redis中
			data, _ := json.Marshal(kpkUser)
			conn.Set(userKey, data, time.Second * 86400)

			return kpkUser, nil
		}
	} else if err != nil {
		return nil, err
	}
	// 从缓存中获取
	kpkUser := new(KpkUser)
	json.Unmarshal([]byte(kpkUserStr), kpkUser)

	return kpkUser, nil
}

// ClearKpkUserCache - 清除用户缓存
func ClearKpkUserCache(UID int64) error {
	conn := conf.Redis.RedisConn
	userKey := fmt.Sprintf("KpkUser_%v", UID)
	return conn.Del(userKey).Err()
}

var rankScoreMap = []struct{
	min  int         // 最小值
	max  int         // 最大值
	rank string      // 名称
}{
	{0, 29, "沧海遗珠"},
	{30, 99, "伴读书童"},
	{100, 199, "金牌书童"},
	{200, 499, "天之骄子"},
	{500, 999, "荣登进士"},
	{1000, 1999, "学富五车"},
	{2000, 4999, "翰林学士"},
	{5000, 9999, "成就非凡"},
	{10000, 19999, "隐退鸿儒"},
	{20000, 0, "巅峰至圣"},
}

// GetRankByScore - 通过积分计算出对应的等级
func GetRankByScore(score64 int64) string {
	score := int(score64)
	if score < 0 {
		score = 0
	}
	rank := ""
	for _, md := range rankScoreMap {
		if score >= md.min && (score <= md.max || md.max == 0) {
			rank = md.rank
			break
		}
	}
	return rank
}

