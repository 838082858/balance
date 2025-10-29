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

func GetBalance(c *gin.Context) {
	var req model.GetBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}

	//查找账户
	var (
		balanceAccountId = req.Id
	)
	resp, err, rows := service.Get(balanceAccountId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//报错，非找不到
		log.Println("server error! ", err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
	} else if rows == 0 {
		//没有对应的数据
		log.Println("search user fail! There is no such data. ", err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "search user fail! There is no such data!", nil))
		return
	}

	log.Printf("search user success!Id:%d, balance:%d, Rows Affected;%d, %+v\n", resp.BalanceAccountId, resp.Balance, rows, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "search user success", resp))

}
