// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
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
		id:         id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			// 对h.clients这个map进行写入操作时加锁防止并发
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			// 同理进行加锁
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Lock()

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
