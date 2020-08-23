package ws

import "fmt"

type methodFunc func(*Client, map[string]interface{}) error

var (
	methodMap = map[string]methodFunc{
		"start":        onStart,        // 开始游戏
		"answer":       onAnswer,       // 答题
		"nextQuestion": onNextQuestion, // 下一题
		"throwProps":   onThrowProps,   // 扔道具，待实现
	}
)

// 业务处理派遣
func dispatch(client *Client, message map[string]interface{}) error {
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
func onStart(client *Client, params map[string]interface{}) error {

	return nil
}

// onAnswer - 答题动作
func onAnswer(client *Client, params map[string]interface{}) error {

	return nil
}

// onNextQuestion -- 下一题
func onNextQuestion(client *Client, params map[string]interface{}) error {

	return nil
}

// onThrowProps -- 扔道具
func onThrowProps(client *Client, params map[string]interface{}) error {

	return nil
}