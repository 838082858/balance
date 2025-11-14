package controller

import (
	"context"
	"errors"
	"http-demo/model"
	"http-demo/service"
	"http-demo/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TransferBalance(c *gin.Context) {
	ctx := c.Request.Context()
	// 超时释放
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel() // 确保 context 资源在函数结束时被释放

	//校验
	var req model.TransferBalanceReq
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		log.Println(bindErr.Error())
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, bindErr.Error(), nil))
		return
	}
	// 调用Service
	resp, err := service.TransferService(ctx, &req)

	//错误返回
	//没有这个账户
	if errors.Is(err, utils.ErrFromAccountNotFound) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "search fromAccountId fail! There is no such data.", nil))
		return
	}
	if errors.Is(err, utils.ErrToAccountNotFound) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "search toAccountId fail! There is no such data.", nil))
		return
	}
	//余额不足
	if errors.Is(err, utils.ErrInsufficientBalance) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "Insufficient balance.", nil))
		return
	}
	//已有交易订单处理
	if errors.Is(err, utils.ErrDuplicateRequest) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "Transfer fail! This is a repetitive operation.Refresh and retry.", nil))
		return
	}
	//账户操作冲突
	if errors.Is(err, utils.ErrOptimisticLockConflict) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "Operation conflict. Please refresh and try again.", nil))
		return
	}
	//账户操作超时
	if errors.Is(err, utils.ErrLockWaitTimeOut) {
		log.Println(err)
		c.JSON(http.StatusBadRequest, model.NewErrResponse(http.StatusBadRequest, "Operation timeout.Please refresh and try again.", nil))
		return
	}
	//todo context超时
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println(err)
		c.JSON(http.StatusGatewayTimeout, model.NewErrResponse(http.StatusBadRequest, "Operation timeout.Please refresh and try again.", nil))
		return
	}
	//其他错误
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.NewErrResponse(http.StatusInternalServerError, "server error.", nil))
		return
	}

	//转账成功
	log.Printf("transfer success! fromId:%d, toId:%d, Amount:%d,  %+v\n", resp.FromAccountId, resp.ToAccountId, resp.Amount, resp)
	c.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "transfer success!",
		Data:    resp,
	})
}
