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
	UserID         int64  `json:"user_id"`   // 用户ID
	Score          int64  `json:"score"`     // 积分值
	PetID          int64  `json:"pet_id"`    // 宠物ID
	RoadID         int64  `json:"road_id"`   // 跑道ID
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

// GetKpkScore -
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
		PetID: data["pet_id"].(int64),
		RoadID: data["road_id"].(int64),
		UpdateTs: data["update_ts"].(string),
	}, nil
}

// KpkUser -用户信息
type KpkUser struct{
	KpkScore
	UserName    string    `json:"user_name"`
	Sex         int32     `json:"sex"`
	Avatar      string    `json:"avatar"`
	Birthday    string    `json:"birthday"`
	City        string    `json:"city"`
	District    string    `json:"district"`
	RoleID      int32     `json:"role_id"`
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