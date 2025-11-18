package mysql_dao

import (
	"context"
	"http-demo/model/mysql_model"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//	func MysqlStorage() {
//		dsn := "root:838082858@tcp(127.0.0.1:3306)/http?charset=utf8mb4&parseTime=True&loc=Local"
//		var dbErr error
//		DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//		if dbErr != nil {
//			panic(dbErr)
//		}
//	}
var (
	daoInstance *Dao
	once        sync.Once
)

type Dao struct {
	db *gorm.DB
}

func GetDaoSingleton() *Dao {
	once.Do(func() {
		db, err := NewMysqlStorage()
		if err != nil {
			panic(err)
		}
		daoInstance = &Dao{db: db}
	})
	return daoInstance
}

// 数据库的初始化
func NewMysqlStorage() (*gorm.DB, error) {
	dsn := "root:838082858@tcp(127.0.0.1:3306)/http?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Dao) Transaction(fc func(tx *gorm.DB) error) error {
	return d.db.Transaction(fc)
}

// select balance
func (d *Dao) GetBalance(ctx context.Context, obj *mysql_model.Balance, balanceAccountId uint64) error {
	//where里面写的字段是数据库里的。SELECT * FROM `balance` WHERE balance_account_id = 555,Find(搜索出来的结果).
	//查重用First()
	findResult := d.db.WithContext(ctx).Where("balance_account_id = ?", balanceAccountId).First(obj)
	return findResult.Error
}

// create balance
func (d *Dao) CreateBalance(ctx context.Context, obj *mysql_model.Balance) error {
	createResult := d.db.WithContext(ctx).Create(obj)
	return createResult.Error
}

// delete balance
func (d *Dao) DeleteBalance(ctx context.Context, obj *mysql_model.Balance) error {
	deleteResult := d.db.WithContext(ctx).Where("balance_account_id = ?", obj.BalanceAccountId).Delete(&mysql_model.Balance{})
	return deleteResult.Error
}

// update balance
func (d *Dao) UpdateBalance(ctx context.Context, obj *mysql_model.Balance) error {
	updateResult := d.db.WithContext(ctx).Model(&mysql_model.Balance{}).Where("balance_account_id = ?", obj.BalanceAccountId).Update("balance", obj.Balance)
	return updateResult.Error
}

func (d *Dao) UpdateTransferBalance(ctx context.Context, obj *mysql_model.Balance) (int64, error) {
	updateResult := d.db.WithContext(ctx).Model(&mysql_model.Balance{}).Where("balance_account_id = ? and version = ?", obj.BalanceAccountId, obj.Version-1).Updates(map[string]interface{}{"balance": obj.Balance, "version": obj.Version})
	return updateResult.RowsAffected, updateResult.Error
}

func (d *Dao) CreateTransfer(ctx context.Context, obj *mysql_model.Transaction) error {
	createResult := d.db.WithContext(ctx).Create(obj)
	return createResult.Error
}

// GetTransfer todo
func (d *Dao) GetTransfer(ctx context.Context, obj *mysql_model.Transaction) error {
	findResult := d.db.WithContext(ctx).Where("transaction_id = ?", obj.TransactionId).First(obj)
	return findResult.Error
}
