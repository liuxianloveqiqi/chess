package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"user/api/internal/config"
	"user/api/internal/middleware"
	"user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	JWT    rest.Middleware
	Rpc    userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,
		JWT:    middleware.NewJWTMiddleware().Handle,
		Rpc:    userclient.NewUser(zrpc.MustNewClient(c.Rpc)),
	}
}
