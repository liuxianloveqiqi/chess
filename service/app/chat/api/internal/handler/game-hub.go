// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"sync"
	"time"
)

// GameHub maintains the set of active clients and broadcasts messages to the
// clients.
type GameHub struct {
	// room id
	id int64
	// Registered clients.
	clients map[*GameClient]bool
	// 系统消息
	systemBroadcast chan []byte

	// Register requests from the clients.
	register chan *GameClient

	// Unregister requests from clients.
	unregister chan *GameClient

	// 添加心跳检测
	heartBeat *time.Ticker

	// 添加互斥锁锁
	mutex sync.Mutex
}

func NewGameHub(id int64) *GameHub {
	return &GameHub{
		id:              id,
		systemBroadcast: make(chan []byte),
		register:        make(chan *GameClient),
		unregister:      make(chan *GameClient),
		clients:         make(map[*GameClient]bool),
	}
}
func (h *GameHub) Run() {
	for {
		select {
		case client := <-h.register:

			// 检查是否已经有两位用户
			if len(h.clients) == 2 {
				// 如果已经有两位用户，则拒绝新用户连接
				client.send <- []byte("房间已经满了，无法加入")
				client.conn.Close()
			} else if len(h.clients) == 1 && h.clients[client] {
				fmt.Println("该用户已经加入")
			} else {
				// 否则将新用户添加到客户端列表中，并向其他用户发送通知
				h.clients[client] = true
				client.send <- []byte(fmt.Sprintf("一个新用户 id= %v 加入了房间room: %v 号", client.id, h.id))
				fmt.Println("客户端的数量为,", len(h.clients))
			}

		case message := <-h.systemBroadcast:
			for client := range h.clients {
				select {
				case client.send <- message:

				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}
