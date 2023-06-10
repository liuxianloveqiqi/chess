package logic

import (
	"context"

	"chat/api/internal/svc"
	"chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameLogic {
	return &GameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameLogic) Game(req *types.JoinRoomReq) error {
	// todo: add your logic here and delete this line

	return nil
}
