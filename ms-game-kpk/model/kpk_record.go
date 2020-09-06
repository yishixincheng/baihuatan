package model

import (
	"baihuatan/pkg/mysql"
	"baihuatan/pkg/common"
	"log"
	"time"
	"fmt"
)

// KpkRecord 用户PK积分记录
type KpkRecord struct {
	ID           int64    `gorose:"id" json:"id"`   // ID
	UserID       int64    `gorose:"user_id" json:"user_id"`   // 用户ID
	Score        int64    `gorose:"score" json:"score"`     // 用户获得积分
	RoomID      string   `gorose:"room_id" json:"room_id"`  // 房间id
	Ranking      int64    `gorose:"ranking" json:"ranking"`   // 名次
	QuestionCount  int64  `gorose:"question_count" json:"question_count"` // 题目数量
	AnswerCount    int64  `gorose:"answer_count" json:"answer_count"`   // 回答数量
	AnswerCorrectCount int64 `gorose:"answer_correct_count" json:"answer_correct_count"` // 回答正确的数量
	UpdateTs       string  `gorose:"update_ts" json:"update_ts"` // 更新时间
}

// TableName -
func (t KpkRecord) TableName() string {
	return "kpk_record";
}

// KpkRecordModel -
type KpkRecordModel struct {
	topNum     int64    // 获取前100个
}

// NewKpkRecordModel -
func NewKpkRecordModel() *KpkRecordModel {
	return &KpkRecordModel{
		topNum: 100,
	}
}

// BatchAdd -
func (p *KpkRecordModel) BatchAdd(data []map[string]interface{}) bool {
	insertData := []map[string]interface{}{
	}
	for _, item := range data {
		insertData = append(insertData, map[string]interface{}{
			"user_id" : item["userID"],
			"room_id" : item["roomID"],
			"score"   : item["score"],
			"question_count" : item["questionCount"],
			"answer_count" : item["answerCount"],
			"answer_correct_count" : item["answerCorrectCount"],
			"ranking" : item["ranking"],
		})
	}

	conn := mysql.DB()
	_, err := conn.Table(&KpkRecord{}).Data(insertData).Insert()
	
	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}

// GetTopListByDay - 日榜
func (p *KpkRecordModel) GetTopListByDay() ([]KpkRecord, error) {
	return p.getTopListByCondition("update_ts", ">=", common.FormatTime(common.GetZeroTime(time.Now()),"Y-m-d H:i:s"))
}

// GetTopListByWeek - 周榜
func (p *KpkRecordModel) GetTopListByWeek() ([]KpkRecord, error) {
	return p.getTopListByCondition("update_ts", ">=", common.FormatTime(common.GetMondayDate(time.Now()),"Y-m-d H:i:s"))
}

// GetTopListByMonth - 月榜
func (p *KpkRecordModel) GetTopListByMonth() ([]KpkRecord, error) {
	return p.getTopListByCondition("update_ts", ">=", common.FormatTime(common.GetFirstDateOfMonth(time.Now()),"Y-m-d H:i:s"))
}

// getTopListByCondition - 
func (p *KpkRecordModel) getTopListByCondition(where ...interface{}) ([]KpkRecord, error) {
	conn := mysql.DB()
	var kpkRList []KpkRecord
	err := conn.Table(&kpkRList).Fields("user_id, sum(score) as score").Where(where...).Group("user_id").Order("score desc").Limit(int(p.topNum)).Select()

	if err == nil {
		return kpkRList, nil
	}
	log.Printf("Error: %v", err)
	return nil, err
}
