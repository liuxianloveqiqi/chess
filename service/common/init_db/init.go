package init_db

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// gorm初始化
func InitGorm(MysqlDataSourece string) *gorm.DB {
	// 将日志写进kafka
	err := logx.Close()
	if err != nil {
		return nil
	}
	db, err := gorm.Open(mysql.Open(MysqlDataSourece),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix:   "tech_", // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		})
	if err != nil {
		panic("连接mysql数据库失败, error=" + err.Error())
	} else {
		fmt.Println("连接mysql数据库成功")
	}
	return db
}

// redis初始化
func InitRedis(add string) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: add,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logx.Error("连接redis失败, error=" + err.Error())
		panic("连接redis失败, error=" + err.Error())
	}
	fmt.Println("redis连接成功")
	return rdb
}

// nsq product初始化
func InitNsqProduct(addr string) *nsq.Producer {
	// 创建NSQ生产者
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		logx.Error("连接redis失败, error=" + err.Error())
		panic("连接redis失败, error=" + err.Error())
	}
	fmt.Println("连接nsq成功！")
	return producer
}
