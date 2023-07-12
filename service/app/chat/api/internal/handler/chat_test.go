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
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODYyOTQwODgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiQVIifQ.Ayv8foFvRcH2zkWFZgOr6b1gDck5X_MrkrE-IItP8G4 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwic3RhdGUiOiJjYjk4MThkZi0xNGYwLTQ2ZWQtOTAzMi03N2MzYjFmZmMyMzAiLCJleHAiOjE2ODg4ODU5MDgsImlhdCI6MTY4NjI5MzkwOCwiaXNzIjoiUlQifQ.ZG42StV1Sw07FM3vxNAg4wHJhBIioeZieEDR_Ey05ZY")

	// 构建 WebSocket 连接地址
	url := "ws://localhost:4003/room/?roomId=1"

	// 连接到 WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {

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
	url := "ws://localhost:4003/room/?roomId=1"

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
	client := &Client{
		sensitiveWords: []string{"敏感词", "关键词"},
	}

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
	client := &Client{
		sensitiveWords: []string{"傻逼", "关键词"},
	}

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
