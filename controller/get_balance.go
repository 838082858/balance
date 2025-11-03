package controller

import (
	"http-demo/model"
	"http-demo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	// 校验
	var req model.GetBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}

	//查找账户
	resp, err, rows := service.Get(&req)

	// 错误判断
	if err != nil && err.Error() == "search user fail! There is no such data." {
		// 没找到
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "search user fail! There is no such data!", nil))
		return
	} else if err != nil {
		// 其他错误
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	// 查找成功
	log.Printf("search user success!Id:%d, balance:%d, Rows Affected;%d, %+v\n", resp.BalanceAccountId, resp.Balance, rows, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "search user success", resp))

}
