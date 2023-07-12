package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
	"user/common/initdb"
	"user/model"
	"user/rpc/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
	Rdb       *redis.Client
	Mdb       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	coon := sqlx.NewMysql(c.Mysql.DataSource)
	MysqlDb := initdb.InitGorm(c.Mysql.DataSource)
	MysqlDb.AutoMigrate(&model.User{})

	redisDb := initdb.InitRedis(c.RedisClient.Host)
	return &ServiceContext{
		Config:    c,
		Mdb:       MysqlDb,
		UserModel: model.NewUserModel(coon, c.CacheRedis),
		Rdb:       redisDb,
	}
}
