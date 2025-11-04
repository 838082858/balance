package controller

import (
	"http-demo/model"
	"http-demo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteBalance(c *gin.Context) {
	ctx := c.Request.Context()
	// 校验
	var req model.DeleteBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadGateway, bindErr.Error(), nil))
		return
	}

	// 删除
	resp, err := service.DeleteService(ctx, &req)

	// 错误返回
	if err != nil && err.Error() == "delete user fail! There is no such data!" {
		// 查找，没有
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "delete user fail! There is no such data!", nil))
		return
	} else if err != nil {
		// 其他错误
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	// 删除成功
	log.Printf("delete user success!Id:%d, balance:%d, %+v\n", resp.BalanceAccountId, resp.Balance, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "delete user success!", resp))

}
