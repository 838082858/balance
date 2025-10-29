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

func UpdateBalance(c *gin.Context) {
	var req model.UpdateBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}

	//查找
	var (
		balanceAccountId = req.Id
	)
	_, getErr, _ := service.Get(balanceAccountId)
	//findResult := dao.DB.Where("balance_account_id", acid.Id).Find(&userBalance)
	if getErr != nil && !errors.Is(getErr, gorm.ErrRecordNotFound) {
		//查找，失败
		log.Println(getErr.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	} else if getErr != nil {
		//查找，没有
		log.Println("search user fail! There is no such data.")
		c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "update user fail! There is no such data!", nil))
		return
	}

	//更新
	updateErr, rows := service.Update(req, req.Balance)
	if updateErr != nil {
		log.Println("update user fail!", updateErr.Error)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "update user fail!server error!", nil))
		return
	}

	//再查找一次
	exist, newGetErr, _ := service.Get(balanceAccountId)
	//findResult := dao.DB.Where("balance_account_id", acid.Id).Find(&userBalance)
	if newGetErr != nil && !errors.Is(newGetErr, gorm.ErrRecordNotFound) {
		//查找，失败
		log.Println(newGetErr.Error())
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error!", nil))
		return
	} else if newGetErr != nil {
		//查找，没有
		log.Println("search user fail! There is no such data.")
		c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "search user fail! There is no such data!", nil))
		return
	}

	var resp = model.UpdateBalanceResp{
		BalanceAccountId: exist.BalanceAccountId,
		Balance:          exist.Balance,
		CreateTime:       exist.CreateTime,
		UpdateTime:       exist.UpdateTime,
		Currency:         exist.Currency,
	}

	log.Printf("update user success!Id:%d, balance:%d, Rows Affected;%d, %+v\n", resp.BalanceAccountId, resp.Balance, rows, resp)
	c.JSON(http.StatusOK, model.NewErrResponse(http.StatusOK, "update success!", resp))

}
