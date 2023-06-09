package logic

import (
	"chess/service/app/user/model"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chess/service/app/user/rpc/internal/svc"
	"chess/service/app/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoReq) (*user.UserList, error) {
	// todo: add your logic here and delete this line

	user0 := model.User{}
	r := l.svcCtx.Mdb.Where("id = ?", in.UserId).First(&user0)
	if r.RowsAffected == 0 {
		return nil, errors.New("数据库查询错误")
	}
	if r.Error != nil {
		return nil, errors.New(r.Error.Error())
	}
	users := make([]*user.UserInfo, 0)
	user1 := &user.UserInfo{
		UserId:     user0.Id,
		PassWord:   user0.Password,
		User_Nick:  user0.UserNick,
		User_Sex:   user0.UserSex,
		User_Email: user0.UserEmail,
		User_Phone: user0.UserPhone,
		CreateTime: timestamppb.New(user0.CreateTime),
		UpdateTime: timestamppb.New(user0.UpdateTime),
		DeleteTime: timestamppb.New(user0.DeleteTime.Time),
	}
	users = append(users, user1)
	logx.Info("这里是users:   ", users)
	return &user.UserList{
		Users: users,
	}, nil

}
