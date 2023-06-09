// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// room id
	id int64
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	systemBroadcast chan []byte

	// Unregister requests from clients.
	unregister chan *Client

	// 添加心跳检测
	heartBeat *time.Ticker

	// 添加互斥锁锁
	mutex sync.Mutex
}

type Message struct {
	client  *Client
	typeS   string //类型
	content []byte
	created time.Time
}

func NewHub(id int64) *Hub {
	return &Hub{
		id:              id,
		broadcast:       make(chan []byte),
		systemBroadcast: make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			// 检查是否已经有两位用户
			if len(h.clients) == 2 {
				// 如果已经有两位用户，则拒绝新用户连接
				client.send <- []byte("房间已经满了，无法加入")
				client.conn.Close()
			} else {
				// 否则将新用户添加到客户端列表中，并向其他用户发送通知
				h.clients[client] = true
				for c := range h.clients {
					if c != client {
						c.send <- []byte(fmt.Sprintf("一个新用户加入了房间room: %v 号", h.id))
					}
				}
				fmt.Println("客户端的数量为,", len(h.clients))
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			// 同理进行加锁
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:

				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}
