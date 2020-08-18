package model

import (
	conf "baihuatan/pkg/config"
	"baihuatan/pkg/mysql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// KpkQuestion 题目列表
type KpkQuestion struct {
	ID           int64    `json:"id"`        // ID
	Title        string   `json:"title"`     // 标题
	Option1      string   `json:"option_1"`  // 选项1
	Option2      string   `json:"option_2"`  // 选项2
	Option3      string   `json:"option_3"`  // 选项3
	Option4      string   `json:"option_4"`  // 选项4
	
}

// KpkQuestionEx 扩展信息
type KpkQuestionEx struct {
	RightOption  string   `json:"right_option"`         // 正确选项
	AuthorID     int64    `json:"author_id"`
	CateID       int64    `json:"cate_id"`
	UpdateTs       string  `json:"update_ts"` // 更新时间
}

// KpkQuestionAll -
type KpkQuestionAll struct {
	KpkQuestion
	KpkQuestionEx
}

// NewKpkQuestionModel -
func NewKpkQuestionModel() *KpkQuestionModel {
	return &KpkQuestionModel{
	}
}

// KpkQuestionModel -
type KpkQuestionModel struct {
}

// getTableName -
func (p *KpkQuestionModel) getTableName() string {
	return "kpk_question"
}

// Add -
func (p *KpkQuestionModel) Add(kpkQuestion *KpkQuestionAll) (int64, error) {
	conn := mysql.DB()
	return conn.Table(p.getTableName()).Data(map[string]interface{}{
		"title" : kpkQuestion.Title,
		"option_1" : kpkQuestion.Option1,
		"option_2" : kpkQuestion.Option2,
		"option_3" : kpkQuestion.Option3,
		"option_4" : kpkQuestion.Option4,
		"right_option" : kpkQuestion.RightOption,
		"author_id" : kpkQuestion.AuthorID,
		"cate_id" : kpkQuestion.CateID,
	}).Insert()
}

// Edit -
func (p *KpkQuestionModel) Edit(kpkQuestion *KpkQuestionAll) (int64, error) {
	if kpkQuestion.ID == 0 {
		return 0, errors.New("id must exist")
	}
	conn := mysql.DB()
	return conn.Table(p.getTableName()).Data(map[string]interface{}{
		"title" : kpkQuestion.Title,
		"option_1" : kpkQuestion.Option1,
		"option_2" : kpkQuestion.Option2,
		"option_3" : kpkQuestion.Option3,
		"option_4" : kpkQuestion.Option4,
		"right_option" : kpkQuestion.RightOption,
		"author_id" : kpkQuestion.AuthorID,
		"cate_id" : kpkQuestion.CateID,
	}).Where("id", kpkQuestion.ID).Update()
}

// GetQuestionListFromCache - 从缓存中获取答题列表
func (p *KpkQuestionModel) GetQuestionListFromCache(num int64) ([]*KpkQuestion, error) {
	conn := conf.Redis.RedisConn
	questionKey := "QuestionStore"
	len, _ := conn.HLen(questionKey).Result()
	if len == 0 {
		return nil, errors.New("empty question")
	}
	if num > len {
		num = len
	}
	rand.Seed(time.Now().UnixNano()) // 随机种子

	var i int64 = 0
	var rds = []string{}

	C: 
	for i = 0; i < num; i++ {
		rd := rand.Intn(int(len)) // 产生随机数
		for _, v := range rds {
			if v == strconv.Itoa(rd) {
				continue C
			}
		}
		rds = append(rds, strconv.Itoa(rd))
	}

	kpkQJSONList, err := conn.HMGet(questionKey, rds...).Result()

	if err != nil {
		return nil, err
	}
	
	kpkQuestionList := []*KpkQuestion{}

	for _, v := range kpkQJSONList {
		var kpkQuestionNode *KpkQuestion
		json.Unmarshal([]byte(v.(string)), kpkQuestionNode)
		kpkQuestionList = append(kpkQuestionList, kpkQuestionNode)
	}

	return kpkQuestionList, nil
}


