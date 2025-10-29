package main

import (
	"http-demo/controller"
	"http-demo/dao/mysql_dao"
	"log"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	//redis启动
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 密码，如果没有则留空
		DB:       0,                // 使用默认的数据库 0,redis有0~15的数据库，之间的数据是隔离的。
	})
	pong, pingErr := rdb.Ping(ctx).Result()
	if pingErr != nil {
		log.Println("redis connection failed:", pingErr)
		return // 如果连接失败，则退出程序
	}
	log.Println("redis connection success:", pong) //成功连接后，pong的值为PONG

	//mysql启动
	mysql_dao.MysqlStorage()

	//gin启动
	r := gin.Default()
	//--------------------数据库create信息
	r.POST("/CreateBalance", controller.CreateBalance)
	//--------------------数据库select信息
	r.POST("/GetBalance", controller.GetBalance)
	//--------------------数据库update信息
	r.POST("/UpdateBalance", controller.UpdateBalance)
	//--------------------数据库delete信息
	r.POST("/DeleteBalance", controller.DeleteBalance)
	err := r.Run()
	if err != nil {
		panic(err)
	}

}
