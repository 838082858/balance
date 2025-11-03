package service

import (
	"errors"
	"http-demo/dao/mysql_dao"
	"http-demo/model"
	"http-demo/model/mysql_model"

	"gorm.io/gorm"
)

func Get(req *model.GetBalanceReq) (*model.GetBalanceResp, error, int64) {
	//不要用var req *model.Balance，这个只声明没有初始化没有分配内存，是指向nil的空指针，First(req)就会报错。
	balance := mysql_model.Balance{}

	// get
	err, row := mysql_dao.GetBalance(&balance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没查到，findResult.Error == gorm.ErrRecordNotFound
		return nil, errors.New("search user fail! There is no such data."), 0
	} else if err != nil {
		// SQL 执行出错（非“记录不存在”）
		return nil, err, 0
	}
	return &model.GetBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		CreateTime:       balance.CreateTime,
		UpdateTime:       balance.UpdateTime,
		Currency:         balance.Currency,
	}, nil, row

}

func Create(req *model.CreateBalanceReq) (*model.CreateBalanceResp, error, int64) {
	balance := mysql_model.Balance{}

	// get是否存在
	err, _ := mysql_dao.GetBalance(&balance, req.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//SQL 执行出错（非“记录不存在”）
		//err.Error == gorm.ErrRecordNotFound没查到数。
		return nil, err, 0
	} else if err == nil {
		//查到记录
		return nil, errors.New("user existed!"), 0
	}

	// create
	balance = mysql_model.Balance{BalanceAccountId: req.Id, Balance: req.Balance, Currency: req.Currency, Version: req.Version}
	err, row := mysql_dao.CreateBalance(&balance)
	if err != nil {
		return nil, errors.New("user create fail，server error!"), 0
	}
	return &model.CreateBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		Currency:         balance.Currency,
	}, nil, row

	//todo redis在这里写
}

func Delete(req *model.DeleteBalanceReq) (*model.DeleteBalanceResp, error, int64) {
	// get
	balance := mysql_model.Balance{}
	err, _ := mysql_dao.GetBalance(&balance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到数据
		return nil, errors.New("delete user fail! There is no such data!"), 0
	} else if err != nil {
		// 查到失败
		return nil, err, 0
	}

	// delete
	err, row := mysql_dao.DeleteBalance(&balance)
	if err != nil {
		//删除失败
		return nil, err, 0
	}
	return &model.DeleteBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		CreateTime:       balance.CreateTime,
		Currency:         balance.Currency,
	}, nil, row

}

func Update(req *model.UpdateBalanceReq) (*model.UpdateBalanceResp, error, int64) {
	// get
	balance := mysql_model.Balance{}
	err, _ := mysql_dao.GetBalance(&balance, req.Id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到数据
		return nil, errors.New("Update user fail! There is no such data!"), 0
	} else if err != nil {
		// 查到失败
		return nil, err, 0
	}

	// update
	balance = mysql_model.Balance{BalanceAccountId: req.Id, Balance: req.Balance}
	err, row := mysql_dao.UpdateBalance(&balance)
	if err != nil {
		return nil, err, 0
	}

	return &model.UpdateBalanceResp{
		BalanceAccountId: balance.BalanceAccountId,
		Balance:          balance.Balance,
		Currency:         balance.Currency,
	}, nil, row

}
