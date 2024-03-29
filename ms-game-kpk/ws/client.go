package ws

import (
	"baihuatan/ms-game-kpk/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message 消息体
type Message map[string]interface{}

// EncodeMessage -
func EncodeMessage(m Message) ([]byte, error) {
	return json.Marshal(m)
}

// DecodeMessage - 
func DecodeMessage(b []byte) (Message, error) {
	m := Message{}
	err := json.Unmarshal(b, &m)
	return m, err
}

//Indicator 答题过程指标
type Indicator struct {
	cursor      int      // 当前用户答题游标
	right       int      // 答正确多少题
	count       int      // 答了多少题
	pace        int      // 步伐，答对加1，答错减一，最低为0
	isAsk       bool     // 答题完成置真，置真的情况下方可以读取下一题
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	user *model.KpkUser
	// The websocket connection.
	conn *websocket.Conn

	// 几人赛
	personNum int64

	// Room指针
	room   *Room

	// 答题指标
	indicator *Indicator

	// 发送通道
	send chan []byte
}


// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("websocket.IsUnexpectedCloseError:", err)
			}
			log.Println("readpump close")
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 读取消息交给处理中心处理
		messageJSONData := make(Message)
		if err := json.Unmarshal(message, &messageJSONData); err != nil {
			// 无效消息
			log.Println("invalid message parameter: ", string(message))
			break;
		}
		if err := dispatch(c, messageJSONData); err != nil {
			log.Println("onMessage Error", err)
			break;
		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				fmt.Println("writepump close")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("nextWrite:", err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("setdealline:", err)
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(ctx context.Context, roomM *RoomManager, w http.ResponseWriter, r *http.Request) {
	// 获取用户信息
	fmt.Println(r.Header.Get("Bhtuser"));
	joinRoomType := r.URL.Query().Get("joinroomtype") // 房间类型
	fmt.Println("房间类型", joinRoomType)
	user, err := roomM.GetKpkUser(ctx, r.Header.Get("Bhtuser"))
	if err != nil {
		log.Println("getUser:", err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-Websocket-Protocol": []string{r.Header.Get("Sec-WebSocket-Protocol")},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// 构建链接客户端
	client := &Client{
		user: user,
		conn: conn,
		personNum: func() int64 {
			if joinRoomType == "2" {
				return 4
			}
			return 2
		}(),
		indicator: &Indicator{
			count: 0,
			cursor: 0,
			right: 0,
			isAsk: false,
		},
		send: make(chan []byte, 256),
	}
	// 匹配房间
	_, err = roomM.MatchingRoom(client)

	if err != nil {
		// 移除客户端
		fmt.Println(err)
		roomM.RemoveClientFromRoom(client)
	}

	go client.writePump()
	go client.readPump()

	// 用户加入房间广播，不能放在
	userJoinBroadcast(client)
}