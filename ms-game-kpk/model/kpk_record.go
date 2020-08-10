package model

import (
	"baihuatan/pkg/mysql"
	"baihuatan/pkg/common"
	"log"
	"time"
)

// KpkRecord 用户PK积分记录
type KpkRecord struct {
	ID           int64    `json:"id"`   // ID
	UserID       int64    `json:"user_id"`   // 用户ID
	Score        int64    `json:"score"`     // 用户获得积分
	HouseID      string   `json:"house_id"`  // 房间id
	Ranking      int64    `json:"ranking"`   // 名次
	QuestionCount  int64  `json:"question_count"` // 题目数量
	AnswerCount    int64  `json:"answer_count"`   // 回答数量
	AnswerCorrectCount int64 `json:"answer_correct_count"` // 回答正确的数量
	UpdateTs       string  `json:"update_ts"` // 更新时间
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

// Add -
func (p *KpkRecordModel) Add(userID, score, ranking, questionCount, answerCount, answerCorrectCount int64, houseID string) bool {
	conn := mysql.DB()
	_, err := conn.Table(&KpkRecord{}).Data(map[string]interface{}{
		"user_id" : userID,
		"score" : score,
		"house_id" : houseID,
		"ranking" : ranking,
		"question_count" : questionCount,
		"answer_count" : answerCount,
		"answer_correct_count" : answerCorrectCount,
	}).Insert()
	
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
