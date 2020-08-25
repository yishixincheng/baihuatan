package ws

import (
	"baihuatan/api/oauth/pb"
	"baihuatan/ms-game-kpk/model"
	conf "baihuatan/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

var (
	// redis队列键名
	roomQueueKey = "bht_room_queue_key"

)

// RoomIDType 房间ID
type RoomIDType  string

// Room - 房间结构体
type Room struct {
	RoomID         RoomIDType   `json:"room_id"`
	OwnerUID       int64        `json:"owner_uid"`      // 房主id
	ClientList     []*Client    `json:"-"`              // 忽略
	ClientNum      int64        `json:"client_num"`     // 客户端数量
	ClientMaxNum   int64        `json:"client_max_num"` // 最大的客户端数量
	CreateTime     time.Time    `json:"create_time"`    // 创建时间
	QuestionList   []*model.KpkQuestion  `json:"-"`     // 忽略
	Status         int64        `json:"status"`         // 状态，0未开始，1人满待开始，2开始
	QuestionNum    int        `json:"-" `             // 题库数量
	WinPace        int        `json:"-"`              // 胜利的步数
	// Inbound message from the clients.
	broadcasts     chan []byte      
	mutex          sync.Mutex      
}

// listen 监听消息
func (p *Room) listen() {
	for {
		select {
		case message := <-p.broadcasts:
			for _, client := range p.ClientList {
				select {
				case client.send <- message:
				default:
					//不阻塞，执行下一次循环
				}
			}		
		}
	}
}

// sendMsgToClient - 发送消息到客户端
func (p *Room) sendMsgToClient(client *Client, message Message) {
	messageByte, err := EncodeMessage(message)
	if err != nil {
		fmt.Println("invalid message:", err)
		return
	}
	// 塞入通道
	client.send <- messageByte
}

// broadcast -- 广播到所有的客户端
func (p *Room) broadcast(message Message) {
	messageByte, err := EncodeMessage(message)
	if err != nil {
		fmt.Println("invalid message:", err)
		return
	}
	// 塞入通道
	p.broadcasts <- messageByte
}

// RoomManager 房间管理者
type RoomManager struct {
	MaxRoomCount       int  // 支持最多的房间数
	RoomList           map[RoomIDType]*Room
	mutex              sync.Mutex 
}

// CreateRoom - 创建房间
func (p *RoomManager) CreateRoom(client *Client) (*Room, error) {
	roomNum := len(p.RoomList)
	if (roomNum >= p.MaxRoomCount) {
		return nil, errors.New("房间已满，请稍后再试！")
	}

	kpkModel := model.NewKpkQuestionModel()
	qustionList, err := kpkModel.GetQuestionListFromCache(100)
	if err != nil {
		return nil, err
	}
	roomID := RoomIDType(uuid.NewV4().String())
	winPace := 10
	questionNum := len(qustionList)
	if winPace > questionNum {
		winPace = questionNum
	}

	room := &Room{
		RoomID: roomID,
		OwnerUID: client.user.UserID,
		ClientList: []*Client{client},
		ClientNum: 1,
		ClientMaxNum: client.personNum,
		QuestionList: qustionList,
		CreateTime: time.Now(),
		Status: 0,
		WinPace: winPace,
		QuestionNum: questionNum,
	}
	client.room = room
	p.RoomList[room.RoomID] = room
	return room, nil
}

// MatchingRoom - 匹配房间,并获得房间对象
// 从redis队列中获取一个元素如果不存在，则创建
func (p *RoomManager) MatchingRoom(client *Client) (*Room, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// 从redis中读取
	redisConn := conf.Redis.RedisConn
	result, err := redisConn.RPop(roomQueueKey).Result()
	if err != nil || err == redis.Nil {
		// 不存在数据，则创建房间
		room, err := p.CreateRoom(client)
		if err != nil {
			return nil, err
		}
		roomJSONStr, err := json.Marshal(room)
		if err != nil {
			return nil, err
		}

		// 塞入缓存中
		if err = redisConn.LPush(roomQueueKey, roomJSONStr).Err(); err != nil {
			return nil, err
		}
		return room, nil
	}

	// 从redis中获取房间
	var roomObj Room
	err = json.Unmarshal([]byte(result), &roomObj)
	if err != nil {
		return nil, err
	}
	room, ok := p.RoomList[roomObj.RoomID]; 
	if !ok {
		return nil, fmt.Errorf("房间ID：%v 不存在", roomObj.RoomID)
	} 
	// 找到房间
	if room.ClientNum >= room.ClientMaxNum {
		return nil, fmt.Errorf("房间ID：%v 已满", roomObj.RoomID)
	}
	client.room   = room
	room.ClientNum ++
	room.ClientList = append(room.ClientList, client)

	if room.ClientNum >= room.ClientMaxNum {
		// 房间已满
		room.Status = 1   // 人满待开始

		return room, nil
	}

	// 未满队尾插入
	roomJSONStr, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}

	// 塞入缓存中
	if err = redisConn.RPush(roomQueueKey, roomJSONStr).Err(); err != nil {
		return nil, err
	}

	return room, nil
}

// RemoveClientFromRoom - 移除客户端从房间中
func (p *RoomManager) RemoveClientFromRoom(client *Client) {
	roomID := client.room.RoomID
	if roomID == "" {
		return
	}
	room, ok := p.RoomList[roomID]
	if !ok {
		return
	}
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for k, clientNode := range room.ClientList {
		if client.conn == clientNode.conn {
			// 移除一个
			room.ClientList = append(room.ClientList[:k], room.ClientList[k+1:]...)
			room.ClientNum --
			close(client.send) //关闭通道
			client.conn.Close()
			break;
		}
	}

	// 没有客户端则移除Room
	if room.ClientNum == 0 {
		close(room.broadcasts)
		delete(p.RoomList, roomID)
	}
	return
}

// GetKpkUser - 获取登录者信息
func (p *RoomManager) GetKpkUser(ctx context.Context, userToken string) (*model.KpkUser, error) {
	var userDetail = pb.UserDetails{}
	if err := json.Unmarshal([]byte(userToken), &userDetail); err != nil {
		fmt.Println("err:", err);
		return nil, err
	}

	return model.GetKpkUserByUID(ctx, userDetail.UserID)
}
