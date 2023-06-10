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

	// Inbound messages from the clients.
	broadcast chan []byte

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
		broadcast:       make(chan []byte),
		systemBroadcast: make(chan []byte),
		register:        make(chan *GameClient),

		unregister: make(chan *GameClient),
		clients:    make(map[*GameClient]bool),
	}
}
func (h *GameHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			// 检查是否已经有两位用户
			if len(h.clients) == 2 {
				// 如果已经有两位用户，则拒绝新用户连接
				h.systemBroadcast <- []byte("系统：游戏房间已经满了，无法加入")
				client.conn.Close()
			} else {
				// 否则将新用户添加到客户端列表中，并向其他用户发送通知
				h.clients[client] = true
				for c := range h.clients {
					if c != client {
						h.systemBroadcast <- []byte(fmt.Sprintf("系统：一个新用户加入了游戏房间: %v 号", h.id))
						h.systemBroadcast <- []byte("系统：请输入 start 准备开始游戏")
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

		case systemMessage := <-h.systemBroadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- systemMessage:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}
