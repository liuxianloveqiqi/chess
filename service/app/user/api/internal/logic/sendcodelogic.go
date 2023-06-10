package logic

import (
	"context"
	"fmt"
	"user/api/internal/svc"
	"user/api/internal/types"
	"user/common/errorx"
	"user/common/utils"
	"user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendcodeLogic {
	return &SendcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendcodeLogic) Sendcode(req *types.RegisterByPhoneRep) (resp *types.RegisterByPhoneResp, err error) {
	// todo: add your logic here and delete this line
	err = utils.DefaultGetValidParams(l.ctx, req)
	if err != nil {
		return nil, errorx.NewDefaultError(fmt.Sprintf("validate校验错误: %v", err))
	}

	cnt, err := l.svcCtx.Rpc.SendCode(l.ctx, &user.SendCodeReq{UserPhone: req.UserPhone})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}

	return &types.RegisterByPhoneResp{VeCode: cnt.VeCode}, nil

}
