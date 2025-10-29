package service

import (
	"errors"
	"http-demo/dao/mysql_dao"
	"http-demo/model"
	"http-demo/model/mysql_model"

	"gorm.io/gorm"
)

func Get(balanceAccountId uint64) (*model.GetBalanceResp, error, int64) {
	//不要用var obj *model.Balance，这个只声明没有初始化没有分配内存，是指向nil的空指针，First(obj)就会报错。
	var obj = &mysql_model.Balance{}
	//where里面写的字段是数据库里的。SELECT * FROM `balance` WHERE balance_account_id = 555,Find(搜索出来的结果).
	//搜索出来的结果放在userBalance。
	//findResult := dao.DB.Where("balance_account_id = ?", acidStr.Id).Find(&userBalance)
	//查重用First()
	err, findRow := mysql_dao.SelectBalance(obj, balanceAccountId)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//SQL 执行出错（非“记录不存在”）
		//findResult.Error == gorm.ErrRecordNotFound没查到数。
		return nil, err, 0
	} else if err == nil {
		//查到记录
		resp := &model.GetBalanceResp{
			BalanceAccountId: obj.BalanceAccountId,
			Balance:          obj.Balance,
			CreateTime:       obj.CreateTime,
			UpdateTime:       obj.UpdateTime,
			Currency:         obj.Currency,
		}
		return resp, nil, findRow
	}
	//没查到
	return nil, err, 0
}

func Create(obj model.CreateBalanceReq) (error, int64) {
	balance := mysql_model.Balance{BalanceAccountId: obj.Id, Balance: obj.Balance, Currency: obj.Currency, Version: obj.Version}
	err, row := mysql_dao.CreateBalance(&balance)
	if err != nil {
		return err, 0
	}
	return nil, row
}
func Delete(obj model.DeleteBalanceReq) (error, int64) {
	balance := mysql_model.Balance{BalanceAccountId: obj.Id}
	err, row := mysql_dao.DeleteBalance(&balance)
	if err != nil {
		return err, 0
	}
	return nil, row
}

func Update(obj model.UpdateBalanceReq, balance uint64) (error, int64) {
	//还需要改一下名称，以及能不能只传一个参数
	obj2 := mysql_model.Balance{BalanceAccountId: obj.Id}
	err, row := mysql_dao.UpdateBalance(&obj2, balance)
	if err != nil {
		return err, 0
	}
	return nil, row
}
