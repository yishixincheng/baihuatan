package ws

import (
	"baihuatan/ms-game-kpk/model"
	"bytes"
	"context"
	"encoding/json"
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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	user *model.KpkUser
	// The websocket connection.
	conn *websocket.Conn

	// 几人赛
	personNum int64

	// Room指针
	room   *Room

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
			log.Println("websocket read message have err")
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 读取消息交给处理中心处理
		messageJSONData := make(map[string]interface{})
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
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
				return
			}
		}
	}
}

// OnMessage - 处理消息
func (c *Client) OnMessage(message []byte) error {
	data := map[string]interface{}{}
	// 解析json
	if err := json.Unmarshal(message, &data); err != nil {
		// 不是json对象不搭理啥事不做
		return nil
	}
	return dispatch(c, data)
}

// ServeWs handles websocket requests from the peer.
func ServeWs(ctx context.Context, roomM *RoomManager, w http.ResponseWriter, r *http.Request) {
	// 获取用户信息
	user, err := roomM.GetKpkUser(ctx, r.Header.Get("Bhtuser"))
	if err != nil {
		log.Println("getUserFile:", err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 构建链接客户端
	client := &Client{
		user: user,
		conn: conn,
		personNum: 4,
		send: make(chan []byte, 256),
	}
	// 匹配房间
	room, err := roomM.MatchingRoom(client)

	if err != nil {
		// 移除客户端
		roomM.RemoveClientFromRoom(client)
	}

	go room.listen()
	go client.writePump()
	go client.readPump()
}