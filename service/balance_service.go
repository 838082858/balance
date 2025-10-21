package service

import (
	"errors"
	"http-demo/dao"
	"http-demo/model"

	"gorm.io/gorm"
)

func Get(balanceAccountId uint64) (*model.Balance, error, int64) {
	//不要用var obj *model.Balance，这个只声明没有初始化没有分配内存，是指向nil的空指针，First(obj)就会报错。
	var obj = &model.Balance{}
	//where里面写的字段是数据库里的。SELECT * FROM `balance` WHERE balance_account_id = 555,Find(搜索出来的结果).
	//搜索出来的结果放在userBalance。
	//findResult := dao.DB.Where("balance_account_id = ?", acidStr.Id).Find(&userBalance)
	//查重用First()
	err, findRow := dao.SelectBalance(obj, balanceAccountId)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//SQL 执行出错（非“记录不存在”）
		//findResult.Error == gorm.ErrRecordNotFound没查到数。
		return nil, err, 0
	} else if err == nil {
		//查到记录
		return obj, nil, findRow
	}
	//没查到
	return nil, err, 0
}

func Creare(balance *model.Balance) (error, int64) {
	err, row := dao.CreateBalance(balance)
	if err != nil {
		return err, 0
	}
	return nil, row
}
func Delete(balance *model.Balance) (error, int64) {
	err, row := dao.DeleteBalance(balance)
	if err != nil {
		return err, 0
	}
	return nil, row
}

func Update(obj *model.Balance, newBalance uint64) (error, int64) {
	err, row := dao.UpdateBalance(obj, newBalance)
	if err != nil {
		return err, 0
	}
	return nil, row
}
