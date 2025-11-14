package mysql_dao

import (
	"context"
	"http-demo/model/mysql_model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MysqlStorage() {
	dsn := "root:838082858@tcp(127.0.0.1:3306)/http?charset=utf8mb4&parseTime=True&loc=Local"
	var dbErr error
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		panic(dbErr)
	}
}

// select balance
func GetBalance(ctx context.Context, obj *mysql_model.Balance, balanceAccountId uint64) error {
	//where里面写的字段是数据库里的。SELECT * FROM `balance` WHERE balance_account_id = 555,Find(搜索出来的结果).
	//查重用First()
	findResult := DB.WithContext(ctx).Where("balance_account_id = ?", balanceAccountId).First(obj)
	return findResult.Error
}

// create balance
func CreateBalance(ctx context.Context, obj *mysql_model.Balance) error {
	createResult := DB.WithContext(ctx).Create(obj)
	return createResult.Error
}

// delete balance
func DeleteBalance(ctx context.Context, obj *mysql_model.Balance) error {
	deleteResult := DB.WithContext(ctx).Where("balance_account_id = ?", obj.BalanceAccountId).Delete(&mysql_model.Balance{})
	return deleteResult.Error
}

// update balance
func UpdateBalance(ctx context.Context, obj *mysql_model.Balance) error {
	updateResult := DB.WithContext(ctx).Model(&mysql_model.Balance{}).Where("balance_account_id = ?", obj.BalanceAccountId).Update("balance", obj.Balance)
	return updateResult.Error
}

func UpdateTransferBalance(ctx context.Context, d *gorm.DB, obj *mysql_model.Balance) (int64, error) {
	updateResult := d.WithContext(ctx).Model(&mysql_model.Balance{}).Where("balance_account_id = ? and version = ?", obj.BalanceAccountId, obj.Version-1).Updates(map[string]interface{}{"balance": obj.Balance, "version": obj.Version})
	return updateResult.RowsAffected, updateResult.Error
}

func CreateTransfer(ctx context.Context, d *gorm.DB, obj *mysql_model.Transaction) error {
	createResult := d.WithContext(ctx).Create(obj)
	return createResult.Error
}

// GetTransfer todo
func GetTransfer(ctx context.Context, d *gorm.DB, obj *mysql_model.Transaction) error {
	findResult := d.WithContext(ctx).Where("transaction_id = ?", obj.TransactionId).First(obj)
	return findResult.Error
}
