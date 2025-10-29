package controller

import (
	"errors"
	"http-demo/model"
	"http-demo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateBalance(c *gin.Context) {
	var req model.CreateBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}
	var (
		balanceAccountId = req.Id
	)
	exist, err, rows := service.Get(balanceAccountId)
	//用结构体会绕过grom框架，数据库的默认值会生效
	//newBalance := map[string]interface{}{
	//	"BalanceAccountId": req.Id, "Balance": req.Balance}
	//查找，已存在
	if exist != nil {
		log.Println("user existed!")
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "user existed!", nil))
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//查找，失败
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	createErr, rows := service.Create(req)
	//创建失败
	if createErr != nil {
		log.Println(createErr.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "user create fail，server error!", nil))
		return
	}
	resp, err, rows := service.Get(balanceAccountId)
	if err != nil {
		//查找，失败
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	log.Printf("create user success! Id:%d, balance:%d, Rows Affected;%d, %+v\n", resp.BalanceAccountId, resp.Balance, rows, resp)
	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "create user success!",
		Data:    resp,
	})
}
