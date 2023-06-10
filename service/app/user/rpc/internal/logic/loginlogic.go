package logic

import (
	"context"
	"errors"
	"user/common/errorx"
	"user/common/utils"
	"user/model"
	"user/rpc/internal/svc"
	"user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.CommonResp, error) {
	// todo: add your logic here and delete this line
	// todo: add your logic here and delete this line
	user0 := model.User{}
	r := l.svcCtx.Mdb.Where("user_phone = ? or user_email = ?", in.PhoneOrEmail, in.PhoneOrEmail).First(&user0)
	if r.RowsAffected == 0 {
		return nil, errors.New("手机号或者邮箱错误")
	}
	if r.Error != nil {
		return nil, errorx.NewDefaultError(r.Error.Error())
	}

	if !utils.ValidMd5Password(in.PassWord, "liuxian", user0.Password) {
		return nil, errors.New("登陆密码错误")
	}
	return &user.CommonResp{
		UserId: user0.Id,
	}, nil
	return &user.CommonResp{}, nil
}
