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

func DeleteBalance(c *gin.Context) {
	var req model.DeleteBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadGateway, bindErr.Error(), nil))
		return
	}

	//查找对应数据
	var (
		balanceAccountId = req.Id
	)
	exist, getErr, _ := service.Get(balanceAccountId)
	//findResult := dao.DB.Where("balance_account_id", req.Id).Find(&userBalance)
	if getErr != nil && !errors.Is(getErr, gorm.ErrRecordNotFound) {
		//查找，失败
		log.Println(getErr.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	} else if getErr != nil {
		//查找，没有
		log.Println("search user fail! There is no such data.")
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "delete user fail! There is no such data!", nil))
		return
	}

	deleteErr, rows := service.Delete(req)
	if deleteErr != nil {
		//删除，失败
		log.Println(deleteErr)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "delete user fail! server error!", nil))
		return
	}
	var resp = model.DeleteBalanceResp{
		BalanceAccountId: exist.BalanceAccountId,
		Balance:          exist.Balance,
		CreateTime:       exist.CreateTime,
		Currency:         exist.Currency,
	}
	log.Printf("delete user success!Id:%d, balance:%d, Rows Affected;%d, %+v\n", exist.BalanceAccountId, exist.Balance, rows, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "delete user success!", resp))
}
