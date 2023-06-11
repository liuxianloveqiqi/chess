package handler

import (
	"bytes"
	"chat/api/internal/logic"
	"chat/api/internal/svc"
	"chat/api/internal/types"
	"chat/common/response"

	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"net/http"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Minute

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	//心跳周期
	heartBeatPeriod = 10 * time.Second

	// 最大等待时间
	maxWaitTime = 10 * time.Minute
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var h *http.Header

var Hubs = make(map[int64]*Hub)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id  int64
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// 最后发送消息的时间
	lastMessageTime int64

	// 互斥锁
	mutex sync.Mutex

	// 限制发送消息的计时器
	limitSpeak *time.Ticker

	//敏感词列表
	sensitiveWords []string

	isReady bool // 标识客户端是否已准备好开始游戏
}

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinRoomReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		hub := Hubs[req.RoomId]
		if hub == nil {
			hub = NewHub(req.RoomId)
			go hub.Run()
			Hubs[req.RoomId] = hub
		}
		h = &r.Header
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		userID := r.Context().Value("user_id")
		if userID == nil {
			logx.Error("获取user_id错误")
			err = errors.New("获取user_id错误")
		}
		client := &Client{id: userID.(int64), hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client
		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
		l := logic.NewChatLogic(r.Context(), svcCtx)
		err = l.Chat(&req)
		response.Response(w, nil, err) //②
	}
}
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.hub.heartBeat = time.NewTicker(heartBeatPeriod)
	c.limitSpeak = time.NewTicker(3 * time.Second)
	lastMessageTime := time.Now()

	for {

		select {

		// 每10秒触发
		case <-c.hub.heartBeat.C:
			// 检测心跳，执行重连
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("发送消息错误: %v", err)
				c.reconnect()
				return
			}
		// 限制三秒才能发一次言
		case <-c.limitSpeak.C:
			// 每三秒重置最后发送消息的时间
			lastMessageTime = time.Now()
		default:
			// 检查距离上一次发送消息的时间是否超过了3秒
			if time.Now().Sub(lastMessageTime) < 3*time.Second {
				// 超过了3秒才允许发送
				break
			}
			// 更新最后一次消息时间
			lastMessageTime = time.Now()
			// 处理图片数据
			_, message, err := c.conn.ReadMessage()
			if strings.HasPrefix(string(message), "image/") {
				c.hub.broadcast <- message
			} else {

				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				// 敏感词替换
				c.sensitiveWords = []string{"傻逼", "暴力", "淫秽"}
				if c.containsSensitiveWords(message) {
					message = c.replaceSensitiveWords(message)
					fmt.Println("开始替换为：" + string(message))
				}
			}

			message = []byte(fmt.Sprintf("userid = %d的用户说：%s", c.id, string(message)))
			for client := range c.hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.hub.clients, client)
				}
			}
			// 更新最后一次消息时间,使用原子操作更新最后发送消息时间防止并发出错
			atomic.StoreInt64(&c.lastMessageTime, time.Now().Unix())
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
			c.mutex.Lock()
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.mutex.Unlock()
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 对 send 进行写入时加锁
			c.mutex.Lock()
			w.Write(message)
			c.mutex.Unlock()
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
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

// 掉线重连
func (c *Client) reconnect() {
	for {
		time.Sleep(3 * time.Second)
		c.mutex.Lock() // 加锁
		waitTime := time.Now().Unix() - c.lastMessageTime
		c.mutex.Unlock() // 解锁
		header := http.Header{}
		header.Set("Authorization", h.Get("Authorization"))

		if waitTime > int64(maxWaitTime.Seconds()) {
			logx.Error("开始重新连接")
			// 开始连接
			// 将 header 作为参数传递给 Dial 方法
			conn, _, err := websocket.DefaultDialer.Dial(
				fmt.Sprintf("ws://localhost:4003/room/?roomId=%v", strconv.Itoa(int(c.hub.id))),
				header,
			)
			if err != nil {
				logx.Error("重连失败: ", err)
			}
			c.conn = conn
			c.send = make(chan []byte, 256)
			log.Println("重连成功")
			go c.writePump()
			go c.readPump()
			// 成功跳调出循环
			break

		}
	}
}

// 检测消息是否含有敏感词
func (c *Client) containsSensitiveWords(message []byte) bool {
	if len(c.sensitiveWords) == 0 {
		return false
	}
	for _, word := range c.sensitiveWords {
		if strings.Contains(string(message), word) {
			return true
		}
	}
	return false
}

// 替换消息中的敏感词为*
func (c *Client) replaceSensitiveWords(message []byte) []byte {
	for _, word := range c.sensitiveWords {
		message = bytes.Replace(message, []byte(word), []byte(strings.Repeat("*", len(word))), -1)
	}
	return message
}
