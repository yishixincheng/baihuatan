package model

import (
	"baihuatan/pkg/common"
	conf "baihuatan/pkg/config"
	"baihuatan/pkg/mysql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
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
	RightOption  string    `json:"right_option"`         // 正确选项
	Annotation   string    `json:"annotation"`             // 注释
	AuthorID     int64     `json:"author_id"`
	CateID       int64     `json:"cate_id"`
	UpdateTs     time.Time  `json:"update_ts"` // 更新时间
}

// KpkQuestionAll -
type KpkQuestionAll struct {
	ID           int64      `gorose:"id" json:"id"`        // ID
	Title        string     `gorose:"title" json:"title"`     // 标题
	Option1      string     `gorose:"option_1" json:"option_1"`  // 选项1
	Option2      string     `gorose:"option_2" json:"option_2"`  // 选项2
	Option3      string     `gorose:"option_3" json:"option_3"`  // 选项3
	Option4      string     `gorose:"option_4" json:"option_4"`  // 选项4
	RightOption  string     `gorose:"right_option" json:"right_option"`         // 正确选项
	Annotation   string     `gorose:"annotation" json:"annotation"`             // 注释
	AuthorID     int64      `gorose:"author_id" json:"author_id"`
	CateID       int64      `gorose:"cate_id" json:"cate_id"`
	UpdateTs     time.Time  `gorose:"update_ts" json:"update_ts"` // 更新时间
}

// TableName -
func (p *KpkQuestionAll) TableName() string {
	return "kpk_question"
}

// NewKpkQuestionModel -
func NewKpkQuestionModel() *KpkQuestionModel {
	return &KpkQuestionModel{
	}
}

// KpkQuestionModel -
type KpkQuestionModel struct {
	sync.RWMutex
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

// Del -
func (p *KpkQuestionModel) Del(ID int64) (int64, error) {
	conn := mysql.DB()
	return conn.Table(p.getTableName()).Where("id", ID).Delete()
}


// GetQuestionList - 获取问题列表用于后台调用,直接从数据库读取
func (p *KpkQuestionModel) GetQuestionList(keyword string, page, pageSize int64) ([]*KpkQuestionAll, error) {
	// 从数据库中获取
	conn := mysql.DB()

	var questionList = []*KpkQuestionAll{}

	err := conn.Table(&questionList).Order("update_ts desc").Offset(int(page)).Limit(int(pageSize)).Select()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return questionList, nil
}

var (
	questionKey     = "QuestionStore"
	questionByIDKey = "QuestionStoreByID" // 通过id映射
)

// GetQuestionListFromCache - 从缓存中获取答题列表
func (p *KpkQuestionModel) GetQuestionListFromCache(num int64) ([]*KpkQuestion, error) {
	p.RLock()
	defer p.RUnlock()
	conn := conf.Redis.RedisConn
	len, _ := conn.HLen(questionKey).Result()
	if len == 0 {
		return nil, errors.New("empty question")
	}
	if num > len {
		num = len
	}

	var rds = []interface{}{}
	for i := 0; i < int(num); i++ {
		rds = append(rds, strconv.Itoa(i))
	}
	rRds := []string{}
	rds = common.UpsetSlice(rds)

	for _, v := range rds {
		rRds = append(rRds, v.(string))
	}
	
	kpkQJSONList, err := conn.HMGet(questionKey, rRds...).Result()

	if err != nil {
		return nil, err
	}
	
	kpkQuestionList := []*KpkQuestion{}

	for _, v := range kpkQJSONList {
		var kpkQuestionNode = &KpkQuestion{}
		json.Unmarshal([]byte(v.(string)), kpkQuestionNode)
		kpkQuestionList = append(kpkQuestionList, kpkQuestionNode)
	}

	return kpkQuestionList, nil
}

// GetQustionFromCache - 从缓存中获取某题的结果
func (p *KpkQuestionModel) GetQustionFromCache(ID int64) (*KpkQuestionAll, error) {
	redisConn := conf.Redis.RedisConn
	kpkQuestionStr, err := redisConn.HGet(questionByIDKey, strconv.FormatInt(ID, 10)).Result()

	if err == redis.Nil {
		// 不存在， 则读取数据库
		conn := mysql.DB()
		result, err := conn.Table(p.getTableName()).Where("id", ID).First()
		
		if err != nil {
			return nil, err
		}

		return  &KpkQuestionAll{
				ID: result["id"].(int64),
				RightOption: result["right_option"].(string),
				Annotation:  result["annotation"].(string),
				CateID:      result["cate_id"].(int64),
				UpdateTs:    result["update_ts"].(time.Time),
		}, nil
	}
	kpkQuestionAll := &KpkQuestionAll{}

	if err := json.Unmarshal([]byte(kpkQuestionStr), kpkQuestionAll); err != nil {
		return nil, err
	}
	return kpkQuestionAll, nil
}

// AutoFetchQuestionsToCache - 自动获取题目到缓存中，启动一个goroutine
func (p *KpkQuestionModel) AutoFetchQuestionsToCache(num int64) {
	p.Lock()
	defer p.Unlock()

	// 从数据库中获取
	conn := mysql.DB()

	var questionList = []KpkQuestionAll{}

	err := conn.Table(&questionList).Order("rand()").Limit(int(num)).Select()

	if err != nil {
		fmt.Println(err)
		return
	}

	// 保存到redis
	redisConn := conf.Redis.RedisConn

	questionV := make(map[string]interface{})
	questionByIDV := make(map[string]interface{})

	for i, v := range questionList {
		 jv, _ := json.Marshal(v)
		 jvs := string(jv)
		 questionV[strconv.Itoa(i)] = jvs
		 questionByIDV[strconv.FormatInt(v.ID,10)] = jvs
	}

	// 设置到redis中 
	redisConn.HMSet(questionKey, questionV)
	redisConn.HMSet(questionByIDKey, questionByIDV)

	return
}


