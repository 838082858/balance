package dao

import (
	"http-demo/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MysqlStorage() {
	dsn := "root:838082858@tcp(127.0.0.1:3306)/http?charset=utf8mb4&parseTime=True&loc=Local"
	var dbErr error
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("Database connection failed：%v", dbErr)
	}
}

// select balance
func SelectBalance(obj *model.Balance, balanceAccountId uint64) (error, int64) {
	findResult := DB.Where("balance_account_id = ?", balanceAccountId).First(obj)
	return findResult.Error, findResult.RowsAffected
}

// create balance
func CreateBalance(obj *model.Balance) (error, int64) {
	createResult := DB.Create(obj)
	return createResult.Error, createResult.RowsAffected
}

// delete balance
func DeleteBalance(obj *model.Balance) (error, int64) {
	deleteResult := DB.Delete(obj)
	return deleteResult.Error, deleteResult.RowsAffected
}

// update balance
func UpdateBalance(obj *model.Balance, balance uint64) (error, int64) {
	updateResult := DB.Model(obj).Update("balance", balance)
	return updateResult.Error, updateResult.RowsAffected
}
