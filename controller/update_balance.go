package controller

import (
	"http-demo/model"
	"http-demo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateBalance(c *gin.Context) {
	// 校验
	var req model.UpdateBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}

	// 查找
	resp, err, row := service.Update(&req)

	// 错误返回
	if err != nil && err.Error() == "Update user fail! There is no such data!" {
		// 查找，没有
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "Update user fail! There is no such data!", nil))
		return
	} else if err != nil {
		// 其他错误
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	}

	// 修改成功
	log.Printf("update user success!Id:%d, balance:%d, Rows Affected;%d, %+v\n", resp.BalanceAccountId, resp.Balance, row, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "Update success!", resp))

}
