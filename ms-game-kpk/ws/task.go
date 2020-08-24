package ws

import (
	"baihuatan/ms-game-kpk/model"
	"fmt"
	"log"

	"golang.org/x/text/cases"
)

const (
	errNo = 0                   // 正确
	errInvalidUser      = 1     // 非法用户
	errInvalidParameter = 2     // 无效参数
	errDataLost         = 3     // 数据丢失
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
			"err" : errInvalidUser,
			"msg" : "非房主无法开始PK",
		})
		return nil
	}
	// 开始游戏，发送题目
	client.room.Status = 2 //开始
	client.room.broadcast(Message{
		"err" : errNo, 
		"msg" : "开始PK",
		"question": *client.room.QuestionList[client.indicator.cursor],
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
			"err" : errDataLost,
			"msg" : "quesion result is lost",
		})
		return nil
	}

	// 答题是否正确
	if optionToChoice(result.RightOption) == choice {
		// 正确
		
	}


	return nil
}

// onNextQuestion -- 下一题
func onNextQuestion(client *Client, message Message) error {

	return nil
}

// onThrowProps -- 扔道具
func onThrowProps(client *Client, message Message) error {

	return nil
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