package main

import (
	"http-demo/controller"
	"http-demo/dao/mysql_dao"
	"http-demo/dao/redis_dao"

	"github.com/gin-gonic/gin"
)

func main() {
	//redis启动
	redis_dao.RedisStorage()

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
