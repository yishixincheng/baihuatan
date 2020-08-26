package ws

import (
	"baihuatan/ms-game-kpk/model"
	"fmt"
	"log"
	"sort"
)

const (
	errNo = 0                   // 正确
	errInvalidUser      = 1     // 非法用户
	errInvalidParameter = 2     // 无效参数
	errDataLost         = 3     // 数据丢失
	errOther            = 9    // 其他错误
)

type methodFunc func(*Client, Message) error

var (
	methodMap = map[string]methodFunc{
		"start":        onStart,        // 开始游戏
		"answer":       onAnswer,       // 答题
		"nextQuestion": onNextQuestion, // 下一题
		"throwProps":   onThrowProps,   // 扔道具，待实现
	}
)

// 业务处理派遣
func dispatch(client *Client, message Message) error {
	method := message["method"].(string)

	if method == "" {
		return fmt.Errorf("method is null")
	}
	handle, ok := methodMap[method]
	if !ok {
		return fmt.Errorf("method %v is not exist", method)
	}

	return handle(client, message)
}

// onStart - 开始
func onStart(client *Client, params Message) error {
	if client.room.OwnerUID != client.user.UserID {
		// 不是房主无法触发开始条件，发送消息
		log.Println("非房主无法开始PK")
		client.room.sendMsgToClient(client, Message{
			"method" : "error",
			"err" : errInvalidUser,
			"msg" : "非房主无法开始PK",
		})
		return nil
	}
	// 开始游戏，发送题目
	client.room.Status = 2 //开始
	client.room.broadcast(Message{
		"method" : "start", 
		"err" :  errNo,
		"msg" : "开始PK",
		"question": *client.room.QuestionList[client.indicator.cursor],
		"number" : client.indicator.cursor + 1, 
	})

	return nil
}

// onAnswer - 答题动作
func onAnswer(client *Client, message Message) error {
	// 答题
	cursor := message["cursor"].(int)     // 当前答题索引
	choice := message["choice"].(string)  // 选择 A|B|C|D

	if cursor > client.indicator.cursor || cursor > len(client.room.QuestionList) - 1 {
		client.room.sendMsgToClient(client, Message{
			"method" : "error",
			"err" : errInvalidParameter,
			"msg" : "cursor is invalid",
		})
		return nil
	}

	question := client.room.QuestionList[cursor]
	kpkModel := model.NewKpkQuestionModel()
	result, err := kpkModel.GetQustionFromCache(question.ID)

	if err != nil {
		client.room.sendMsgToClient(client, Message{
			"method": "error",
			"err" : errDataLost,
			"msg" : "quesion result is lost",
		})
		return nil
	}

	rightChoice := optionToChoice(result.RightOption)
	// 正确
	client.room.sendMsgToClient(client, Message{
		"method": "answerResult",
		"choice"      : choice,
		"cursor"      : cursor,
		"rightChoice" : rightChoice,
		"annotation"  : result.Annotation,
	})

	client.indicator.isAsk = true

	// 计算战果
	client.indicator.count ++    // 答题总数
	if rightChoice == choice {
		// 争取答题总数
		client.indicator.right ++
	}

	client.indicator.pace = 2 * client.indicator.right - client.indicator.count
	if client.indicator.pace < 0 {
		// 不能小于0
		client.indicator.pace = 0
	}

	responseFightDynamic(client.room)

	// 计算结果，率先答对10题，或者答完，则游戏结束计算战果
	if (client.indicator.pace >= client.room.WinPace) && (client.indicator.count == client.room.QuestionNum) {
		gameOverSummary(client.room)
	}

	return nil
}

// onNextQuestion -- 下一题
func onNextQuestion(client *Client, message Message) error {
	if client.room.Status != 2 {
		client.room.sendMsgToClient(client, Message{
			"method": "error",
			"err" : errOther,
			"msg" : "PK游戏还未开始",
		})
		return nil
	}
	if !client.indicator.isAsk {
		// 上一题没答复无法读取下一题
		client.room.sendMsgToClient(client, Message{
			"method": "error",
			"err" : errOther,
			"msg" : "请答完题再翻下一题",
		})
		return nil
	}

	// 游标移动
	client.indicator.cursor ++
	client.indicator.isAsk = false

	client.room.broadcast(Message{
		"method" : "newquestion", 
		"err" :  errNo,
		"msg" : "请答题",
		"question": *client.room.QuestionList[client.indicator.cursor],
		"number" : client.indicator.cursor + 1, 
	})
	
	return nil
}

// onThrowProps -- 扔道具
func onThrowProps(client *Client, message Message) error {

	return nil
}

// 发送战斗动态
func responseFightDynamic(room *Room) error {
	message := Message{
		"method" : "fightDynamic",
		"roomID" : room.RoomID,
		"status" : room.Status,
	}

	userList := []Message{
	}

	for _, client := range room.ClientList {
		userList = append(userList, Message{
			"userID" : client.user.UserID,
			"userName" : client.user.UserName,
			"pace"     : client.indicator.pace,
			"right"    : client.indicator.right,
			"count"    : client.indicator.count,
		})
	}
	message["userList"] = userList

	// 广播
	room.broadcast(message)

	return nil
}

// 游戏结束汇总, 计算排名，发送消息，并统计到数据库中
func gameOverSummary(room *Room) {
	room.Status = 3
	message := Message{
		"method" : "gameover",
		"roomID" : room.RoomID,
		"status" : room.Status,
	}

	// 求出1，2，3名，并计算获得积分值
	userList := []Message{
	}
	for _, client := range room.ClientList {
		userList = append(userList, Message{
			"userID" : client.user.UserID,
			"userName" : client.user.UserName,
			"pace"     : client.indicator.pace,
			"right"    : client.indicator.right,
			"count"    : client.indicator.count,
			"core"     : 0,
		})
	}

	// 降序排序
	sort.Slice(userList, func(i, j int ) bool {
		return userList[i]["pace"].(int) > userList[i]["pace"].(int)
	})

	// 计算分值
	if room.ClientNum > 2 {
		userList[0]["core"] = 5
		userList[1]["core"] = 2
	} else {
		userList[0]["core"] = 2
	}
	message["userList"] = userList

	// 广播
	room.broadcast(message)

	insertData := []map[string]interface{}{
	}
	
	for i, user := range userList {
		insertData = append(insertData, map[string]interface{}{
			"userID" : user["userID"],
			"roomID" : room.RoomID,
			"score"  : user["score"],
			"questionCount" : room.QuestionNum,
			"answerCount" : user["count"],
			"answerCorrectCount" : user["right"],
			"ranking" : (i+1),
		})
	}
	// 记录到表中
	kpkRecordModel := model.NewKpkRecordModel()
	kpkRecordModel.BatchAdd(insertData)
	kpkScoreModel := model.NewKpkScoreModel()
	kpkScoreModel.BatchIncScore(insertData)
}

// 用户加入房间通知
func userJoinNotify(client *Client) {
	room := client.room
	message := Message{
		"method" : "userjoin",
		"roomID" : room.RoomID,
		"user" : Message{
			"userID": client.user.UserID,
			"userName": client.user.UserName,
			"roadID": client.user.RoadID,
			"petID": client.user.PetID,
			"avatar": client.user.Avatar,
			"sex": client.user.Sex,
			"score": client.user.Score,
			"rank": client.user.Rank,
		},
	}

	// 广播
	room.broadcast(message)
}

// optionToChoice -
func optionToChoice(option string) string {
	switch option {
	case "1" :
		return "A";
	case "2" :
		return "B";
	case "3" :
		return "C";
	case "4":
		return "D"
	}
	return ""
}

