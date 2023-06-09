package svc

import (
	"chess/service/app/user/model"
	"chess/service/app/user/rpc/internal/config"
	"chess/service/common/init_db"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
	Rdb       *redis.Client
	Mdb       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	coon := sqlx.NewMysql(c.Mysql.DataSource)
	MysqlDb := init_db.InitGorm(c.Mysql.DataSource)
	MysqlDb.AutoMigrate(&model.User{})

	redisDb := init_db.InitRedis(c.RedisClient.Host)
	return &ServiceContext{
		Config:    c,
		Mdb:       MysqlDb,
		UserModel: model.NewUserModel(coon, c.CacheRedis),
		Rdb:       redisDb,
	}
}
