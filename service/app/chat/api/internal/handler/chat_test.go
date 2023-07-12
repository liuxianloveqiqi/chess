package handler

import (
	"bytes"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
	"time"
)

func TestWritePump(t *testing.T) {
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2ODg3OTgzNzEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiQVIifQ.Ovx_cOjk-FqJ6T3jtY_LEsu4NFa-uy3O3nkNYGdekYI eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2OTEzOTAxOTEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiUlQifQ.keO4V1EIpYy75uLPCgyd8QlPgo3kuttt1RrlELRgebo")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		recover()
	}

	defer conn.Close()

	client := &Client{
		id: 123,
		hub: &Hub{
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			heartBeat:  time.NewTicker(heartBeatPeriod),
		},
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.clients[client] = true

	// 模拟发送心跳消息
	go func() {
		<-client.hub.heartBeat.C
		client.conn.WriteMessage(websocket.PingMessage, nil)
	}()

	// 启动 writePump 方法
	go client.writePump()

	// 模拟发送消息到 send 通道
	client.send <- []byte("Hello, world!")

	// 验证是否成功发送消息
	msgType, msgBytes, err := client.conn.ReadMessage()
	if err != nil {
		t.Errorf("发送消息错误：%v", err)
	}
	if msgType != websocket.TextMessage {
		t.Errorf("发送的消息类型不正确。期望值：%d，实际值：%d", websocket.TextMessage, msgType)
	}
	if string(msgBytes) != "Hello, world!" {
		t.Errorf("发送的消息内容不正确。期望值：\"Hello, world!\"，实际值：%s", string(msgBytes))
	}
}

func TestReadPump(t *testing.T) {
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODYyOTQwODgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiQVIifQ.Ayv8foFvRcH2zkWFZgOr6b1gDck5X_MrkrE-IItP8G4 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODg4ODU5MDgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiUlQifQ.ZG42StV1Sw07FM3vxNAg4wHJhBIioeZieEDR_Ey05ZY")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		t.Fatalf("无法连接到 WebSocket 服务器：%v", err)
	}
	defer conn.Close()

	client := &Client{
		id: 123,
		hub: &Hub{
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			heartBeat:  time.NewTicker(heartBeatPeriod),
		},
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.clients[client] = true

	client.readPump()

	// 验证是否成功读取消息
	select {
	case msg := <-client.send:
		expectedMsg := "Hello, world!"
		if string(msg) != expectedMsg {
			t.Errorf("读取的消息不正确。期望值：%s，实际值：%s", expectedMsg, string(msg))
		}
	default:
		t.Error("未成功读取消息")
	}
}
func TestContainsSensitiveWords(t *testing.T) {
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2ODg3OTgzNzEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiQVIifQ.Ovx_cOjk-FqJ6T3jtY_LEsu4NFa-uy3O3nkNYGdekYI eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2OTEzOTAxOTEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiUlQifQ.keO4V1EIpYy75uLPCgyd8QlPgo3kuttt1RrlELRgebo")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		recover()
	}

	defer conn.Close()

	client := &Client{
		id: 123,
		hub: &Hub{
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			heartBeat:  time.NewTicker(heartBeatPeriod),
		},
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.clients[client] = true

	tests := []struct {
		message       []byte
		expectedValue bool
	}{
		{[]byte("这是一条测试信息"), false},
		{[]byte("这句话包含敏感词"), true},
		{[]byte("这个消息中有多个关键词"), true},
	}

	for _, test := range tests {
		result := client.containsSensitiveWords(test.message)
		if result != test.expectedValue {
			t.Errorf("测试结果不符合预期。消息：%s，预期：%v，实际：%v",
				string(test.message), test.expectedValue, result)
		}
	}
}

func TestReplaceSensitiveWords(t *testing.T) {
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2ODg3OTgzNzEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiQVIifQ.Ovx_cOjk-FqJ6T3jtY_LEsu4NFa-uy3O3nkNYGdekYI eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2OTEzOTAxOTEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiUlQifQ.keO4V1EIpYy75uLPCgyd8QlPgo3kuttt1RrlELRgebo")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		recover()
	}

	defer conn.Close()

	client := &Client{
		id: 123,
		hub: &Hub{
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte),
			heartBeat:  time.NewTicker(heartBeatPeriod),
		},
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.clients[client] = true

	tests := []struct {
		message       []byte
		expectedValue []byte
	}{
		{[]byte("这是一条测试信息"), []byte("这是一条测试信息")},
		{[]byte("这句话包含敏感词"), []byte("这句话包含*********")},
		{[]byte("这个消息中有多个关键词"), []byte("这个消息中有多个*********")},
	}

	for _, test := range tests {
		result := client.replaceSensitiveWords(test.message)
		if !bytes.Equal(result, test.expectedValue) {
			t.Errorf("测试结果不符合预期。消息：%s，预期：%s，实际：%s",
				string(test.message), string(test.expectedValue), string(result))
		}
	}
}
func TestAreBothClientsReady(t *testing.T) {
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2ODg3OTgzNzEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiQVIifQ.Ovx_cOjk-FqJ6T3jtY_LEsu4NFa-uy3O3nkNYGdekYI eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NSwic3RhdGUiOiI4ZjcxNWQxNi1mOWJmLTQ2NzEtOTVjYy1mMWQ1ODg4ZDdhMmUiLCJleHAiOjE2OTEzOTAxOTEsImlhdCI6MTY4ODc5ODE5MSwiaXNzIjoiUlQifQ.keO4V1EIpYy75uLPCgyd8QlPgo3kuttt1RrlELRgebo")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		recover()
	}

	defer conn.Close()

	client := &GameClient{
		id: 123,
		hub: &GameHub{
			register:        make(chan *GameClient),
			unregister:      make(chan *GameClient),
			clients:         make(map[*GameClient]bool),
			systemBroadcast: make(chan []byte),
			heartBeat:       time.NewTicker(heartBeatPeriod),
		},
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.clients[client] = true

	// 添加两个客户端，并设置一个客户端为未准备状态
	client1 := &GameClient{isReady: true}
	client2 := &GameClient{isReady: false}
	client.hub.clients[client1] = true
	client.hub.clients[client2] = true

	ready := client.hub.areBothClientsReady()

	if ready {
		t.Error("期望返回值为 false，实际为 true")
	}

	// 设置第二个客户端为准备状态
	client2.isReady = true

	ready = client.hub.areBothClientsReady()

	if !ready {
		t.Error("期望返回值为 true，实际为 false")
	}
}

func TestStringWhiteOrBlack(t *testing.T) {
	white := true
	black := false
	h := &GameHub{
		clients: make(map[*GameClient]bool),
	}
	whiteString := h.stringWhiteOrBlack(white)

	if whiteString != "白方" {
		t.Errorf("期望返回值为 \"白方\"，实际为 %s", whiteString)
	}

	blackString := h.stringWhiteOrBlack(black)

	if blackString != "黑方" {
		t.Errorf("期望返回值为 \"黑方\"，实际为 %s", blackString)
	}
}
func TestRun(t *testing.T) {
	hub := NewHub(1)
	// 创建 HTTP 请求头
	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODYyOTQwODgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiQVIifQ.Ayv8foFvRcH2zkWFZgOr6b1gDck5X_MrkrE-IItP8G4 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODg4ODU5MDgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiUlQifQ.ZG42StV1Sw07FM3vxNAg4wHJhBIioeZieEDR_Ey05ZY")

	// 构建 WebSocket 连接地址
	url := "ws://43.139.195.17:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		t.Fatalf("无法连接到 WebSocket 服务器：%v", err)
	}
	defer conn.Close()
	client1 := &Client{
		conn: conn,
		send: make(chan []byte),
	}
	client2 := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	hub.register <- client1
	hub.register <- client2

	time.Sleep(time.Millisecond) // 等待 goroutine 执行注册操作

	if len(hub.clients) != 2 {
		t.Errorf("期望客户端数量为 2，实际为 %d", len(hub.clients))
	}

	message := []byte("Test message")
	hub.broadcast <- message

	time.Sleep(time.Millisecond) // 等待 goroutine 执行广播操作

	for client := range hub.clients {
		receivedMessage := <-client.send
		if string(receivedMessage) != string(message) {
			t.Errorf("期望接收到的消息为 \"%s\"，实际为 \"%s\"", string(message), string(receivedMessage))
		}
	}

	// 模拟一个客户端断开连接
	hub.unregister <- client1

	time.Sleep(time.Millisecond) // 等待 goroutine 执行注销操作

	if len(hub.clients) != 1 {
		t.Errorf("期望客户端数量为 1，实际为 %d", len(hub.clients))
	}
}

func TestNewHub(t *testing.T) {
	id := int64(1)
	hub := NewHub(id)

	if hub.id != id {
		t.Errorf("期望 Hub 的 id 为 %d，实际为 %d", id, hub.id)
	}

	if hub.broadcast == nil {
		t.Error("Hub 的 broadcast 通道未初始化")
	}

	if hub.register == nil {
		t.Error("Hub 的 register 通道未初始化")
	}

	if hub.unregister == nil {
		t.Error("Hub 的 unregister 通道未初始化")
	}

	if hub.clients == nil {
		t.Error("Hub 的 clients 列表未初始化")
	}
}
