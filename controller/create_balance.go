package controller

import (
	"http-demo/model"
	"http-demo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBalance(c *gin.Context) {
	ctx := c.Request.Context()
	//校验
	var req model.CreateBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}

	// 创建
	resp, err := service.CreateService(ctx, &req)

	//错误返回
	if err != nil && err.Error() == "user existed!" {
		//查找已存在
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "user existed!", nil))
		return
	} else if err != nil && err.Error() == "user create fail，server error!" {
		//创建失败
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "user create fail，server error!", nil))
		return
	} else if err != nil {
		//其他错误
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	//创建成功
	log.Printf("create user success! Id:%d, balance:%d,  %+v\n", resp.BalanceAccountId, resp.Balance, resp)
	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "create user success!",
		Data:    resp,
	})
}
