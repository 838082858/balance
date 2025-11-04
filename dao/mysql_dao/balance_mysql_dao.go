package mysql_dao

import (
	"http-demo/model/mysql_model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func MysqlStorage() {
	dsn := "root:838082858@tcp(127.0.0.1:3306)/http?charset=utf8mb4&parseTime=True&loc=Local"
	var dbErr error
	db, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		panic(dbErr)
	}
}

// select balance
func GetBalance(obj *mysql_model.Balance, balanceAccountId uint64) error {
	//where里面写的字段是数据库里的。SELECT * FROM `balance` WHERE balance_account_id = 555,Find(搜索出来的结果).
	//查重用First()
	findResult := db.Where("balance_account_id = ?", balanceAccountId).First(obj)
	return findResult.Error
}

// create balance
func CreateBalance(obj *mysql_model.Balance) error {
	createResult := db.Create(obj)
	return createResult.Error
}

// delete balance
func DeleteBalance(obj *mysql_model.Balance) error {
	deleteResult := db.Where("balance_account_id = ?", obj.BalanceAccountId).Delete(&mysql_model.Balance{})
	return deleteResult.Error
}

// update balance
func UpdateBalance(obj *mysql_model.Balance) error {
	updateResult := db.Model(&mysql_model.Balance{}).Where("balance_account_id = ?", obj.BalanceAccountId).Update("balance", obj.Balance)
	return updateResult.Error
}
