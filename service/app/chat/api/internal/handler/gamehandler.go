package handler

import (
	"chat/api/internal/logic"
	"chat/api/internal/svc"
	"chat/api/internal/types"
	"chat/common/response"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"strconv"
	"sync/atomic"

	"log"
	"net/http"
	"sync"
	"time"
)

var h0 *http.Header

var GameHubs = make(map[int64]*GameHub)

var state State

// GameClient is a middleman between the websocket connection and the GameHub.
type GameClient struct {
	id  int64
	hub *GameHub
	// The websocket connection.
	conn *websocket.Conn

	// 最后发送消息的时间
	lastMessageTime int64

	// Buffered channel of outbound messages.
	send chan []byte
	// 互斥锁
	mutex       sync.Mutex
	isReady     bool // 标识客户端是否已准备好开始游戏
	gameStarted bool // 游戏是否开始
	isBout      bool // 判断是否是自己回合
	// 黑白方
	isWhite bool
}

func gameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinRoomReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		hub := GameHubs[req.RoomId]
		if hub == nil {
			hub = NewGameHub(req.RoomId)
			go hub.Run()
			GameHubs[req.RoomId] = hub
		}
		h0 = &r.Header
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		userID := r.Context().Value("user_id")
		fmt.Println("useri id ssss", userID)
		if userID == nil {
			logx.Error("获取user_id错误")
			err = errors.New("获取user_id错误")
		}
		client := &GameClient{id: userID.(int64), hub: hub, conn: conn, send: make(chan []byte, 256), isReady: false, isWhite: false, isBout: false}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
		l := logic.NewGameLogic(r.Context(), svcCtx)
		err = l.Game(&req)
		response.Response(w, nil, err) //②

	}
}

// 掉线重连
func (c *GameClient) reconnect() {
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
				fmt.Sprintf("ws://localhost:4003/room/game?roomId=%v", strconv.Itoa(int(c.hub.id))),
				header,
			)
			if err != nil {
				logx.Error("重连失败: ", err)
			}
			c.mutex.Lock() // 加锁
			c.conn = conn
			c.mutex.Unlock() // 解锁
			c.conn = conn
			c.send = make(chan []byte, 256)
			logx.Info("重连成功")
			go c.writePump()
			go c.readPump()
			// 成功跳调出循环
			break

		}
	}
}
func (c *GameClient) readPump() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.hub.heartBeat = time.NewTicker(heartBeatPeriod)
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		select {
		// 每10秒触发
		case <-c.hub.heartBeat.C:
			// 检测心跳，执行重连
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("发送消息错误: %v", err)
				c.reconnect()
				c.mutex.Unlock() // 解锁
				return
			}

		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			if !c.isReady && string(message) != "start" {
				c.send <- []byte("系统：请输入start开始")
			}
			if string(message) == "start" {
				c.isReady = true
				if c.hub.areBothClientsReady() && !c.gameStarted {

					c.hub.systemBroadcast <- []byte("系统：游戏开始\n")
					// 定黑白 1白0黑

					for client, _ := range c.hub.clients {
						client.isWhite = true
						client.isBout = true
						fmt.Println(client.id, client.isWhite, client.isBout)
						break

					}

					fmt.Println("现在的客户端的数量", len(c.hub.clients))
					fmt.Println(c.hub.clients)
					c.hub.systemBroadcast <- []byte(fmt.Sprintf("user id为%v的用户为%v\n", c.id, c.hub.stringWhiteOrBlack(c.isWhite)))
					c.gameStarted = true
					state = c.NewInitialBoard(c.isWhite)
					c.hub.systemBroadcast <- []byte(state.board.String())
				}
			} else if !c.hub.areBothClientsReady() && c.isReady {
				c.send <- []byte("系统：请等待对方开始")
			} else {
				fmt.Println("nnnnnnnnnn")
				userMessage := []byte(fmt.Sprintf("userid = %d的用户操作：%s", c.id, string(message)))
				c.hub.systemBroadcast <- userMessage

				// 处理用户输入的消息并获取chess系统的响应消息
				input := c.processInput(string(message))

				// 将引擎的响应消息发送给客户端
				c.hub.systemBroadcast <- input

				// 更新最后一次消息时间,使用原子操作更新最后发送消息时间防止并发出错
				atomic.StoreInt64(&c.lastMessageTime, time.Now().Unix())
			}

		}
	}
}

func (c *GameClient) writePump() {
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
				// The hub closed the channel.
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

func (h *GameHub) areBothClientsReady() bool {
	// 遍历所有客户端，检查是否都已准备好开始游戏
	for client := range h.clients {
		if !client.isReady || len(h.clients) < 2 {
			return false
		}
	}
	return true
}

func (h *GameHub) stringWhiteOrBlack(is bool) string {
	return map[bool]string{true: "白方", false: "黑方"}[is]
}
