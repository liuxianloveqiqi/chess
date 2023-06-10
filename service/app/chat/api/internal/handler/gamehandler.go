package handler

import (
	"chat/api/internal/logic"
	"chat/api/internal/svc"
	"chat/api/internal/types"
	"chat/common/response"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"net/http"
	"sync"
	"time"
)

type GameClient struct {
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

func gameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinRoomReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
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
		l := logic.NewGameLogic(r.Context(), svcCtx)
		err = l.Game(&req)
		response.Response(w, nil, err) //②

	}
}
