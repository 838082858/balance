package service

import (
	"context"
	"errors"
	"http-demo/dao/mysql_dao"
	"http-demo/dao/redis_dao"
	"http-demo/model"
	"http-demo/model/mysql_model"
	"http-demo/model/redis_model"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetService(ctx context.Context, req *model.GetBalanceReq) (*model.GetBalanceResp, error) {
	// redis get
	redisBalance, err := redis_dao.GetBalanceCache(ctx, req.Id)
	if redisBalance != nil && err == nil {
		log.Println("redis get success!")
		return &model.GetBalanceResp{
			BalanceAccountId: redisBalance.BalanceAccountId,
			Balance:          redisBalance.Balance,
			CreateTime:       redisBalance.CreateTime,
			UpdateTime:       redisBalance.UpdateTime,
			Currency:         redisBalance.Currency,
		}, nil
	} else if errors.Is(err, redis.Nil) {
		log.Println(err)
	}

	// sql get
	//不要用var req *model.Balance，这个只声明没有初始化没有分配内存，是指向nil的空指针，First(req)就会报错。
	mysqlBalance := mysql_model.Balance{}
	err = mysql_dao.GetBalance(ctx, &mysqlBalance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没查到，findResult.Error == gorm.ErrRecordNotFound
		log.Println(err)
		return nil, errors.New("search user fail! There is no such data.")
	} else if err != nil {
		// SQL 执行出错（非“记录不存在”）
		log.Println(err)
		return nil, err
	}

	// redis set
	redisBalance = &redis_model.Balance{
		BalanceAccountId: mysqlBalance.BalanceAccountId,
		Balance:          mysqlBalance.Balance,
		CreateTime:       mysqlBalance.CreateTime,
		UpdateTime:       mysqlBalance.UpdateTime,
		Currency:         mysqlBalance.Currency,
	}
	err = redis_dao.SetBalanceCache(ctx, req.Id, redisBalance)
	if err != nil {
		log.Println(err)
	}
	log.Println("mysql get success!")

	// return
	return &model.GetBalanceResp{
		BalanceAccountId: mysqlBalance.BalanceAccountId,
		Balance:          mysqlBalance.Balance,
		CreateTime:       mysqlBalance.CreateTime,
		UpdateTime:       mysqlBalance.UpdateTime,
		Currency:         mysqlBalance.Currency,
	}, nil

}

func CreateService(ctx context.Context, req *model.CreateBalanceReq) (*model.CreateBalanceResp, error) {

	balance := mysql_model.Balance{}

	// get是否存在
	err := mysql_dao.GetBalance(ctx, &balance, req.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//SQL 执行出错（非“记录不存在”）
		//err.Error == gorm.ErrRecordNotFound没查到数。
		return nil, err
	} else if err == nil {
		//查到记录
		return nil, errors.New("user existed!")
	}

	// create
	balance = mysql_model.Balance{BalanceAccountId: req.Id, Balance: req.Balance, Currency: req.Currency, Version: req.Version}
	err = mysql_dao.CreateBalance(ctx, &balance)
	if err != nil {
		return nil, errors.New("user create fail，server error!")
	}
	return &model.CreateBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		Currency:         balance.Currency,
	}, nil

}

func DeleteService(ctx context.Context, req *model.DeleteBalanceReq) (*model.DeleteBalanceResp, error) {
	// sql get
	balance := mysql_model.Balance{}
	err := mysql_dao.GetBalance(ctx, &balance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到数据
		log.Println(err)
		return nil, errors.New("delete user fail! There is no such data!")
	} else if err != nil {
		// 查到失败
		log.Println(err)
		return nil, err
	}
	log.Println("mysql get success!")

	// redis delete
	err = redis_dao.DeleteBalanceCache(ctx, req.Id)
	if err != nil {
		log.Println(err)
	}

	// sql delete
	err = mysql_dao.DeleteBalance(ctx, &balance)
	if err != nil {
		//删除失败
		log.Println(err)
		return nil, err
	}
	log.Println("mysql delete success!")
	return &model.DeleteBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		CreateTime:       balance.CreateTime,
		Currency:         balance.Currency,
	}, nil

}

func UpdateService(ctx context.Context, req *model.UpdateBalanceReq) (*model.UpdateBalanceResp, error) {
	// sql get
	balance := mysql_model.Balance{}
	err := mysql_dao.GetBalance(ctx, &balance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到数据
		log.Println(err)
		return nil, errors.New("Update user fail! There is no such data!")
	} else if err != nil {
		// 查到失败
		log.Println(err)
		return nil, err
	}
	log.Println("mysql get success!")

	// redis delete
	err = redis_dao.DeleteBalanceCache(ctx, req.Id)
	if err != nil {
		log.Println(err)
	}

	// sql update
	balance = mysql_model.Balance{BalanceAccountId: req.Id, Balance: req.Balance}
	err = mysql_dao.UpdateBalance(ctx, &balance)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("mysql update success!")
	return &model.UpdateBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
	}, nil

}
