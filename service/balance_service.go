package service

import (
	"context"
	"errors"
	"http-demo/dao/mysql_dao"
	"http-demo/dao/redis_dao"
	"http-demo/model"
	"http-demo/model/mysql_model"
	"http-demo/model/redis_model"
	"http-demo/utils"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
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
		return nil, errors.New("search user fail! There is no such data")
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
		//查询时发生数据库错误（不是“没查到记录”）
		return nil, err
	} else if err == nil {
		//查到记录
		return nil, utils.ErrAccountFound
	}

	// create
	balance = mysql_model.Balance{BalanceAccountId: req.Id, Balance: req.Balance, Currency: req.Currency}
	err = mysql_dao.CreateBalance(ctx, &balance)
	if err != nil {
		return nil, errors.New("user create fail，server error")
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
		return nil, errors.New("delete user fail! There is no such data")
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
		return nil, errors.New("update user fail! There is no such data")
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

func TransferService(ctx context.Context, req *model.TransferBalanceReq) (*model.TransferBalanceResp, error) {
	fromAccount := mysql_model.Balance{}
	toAccount := mysql_model.Balance{}
	transfer := mysql_model.Transaction{
		TransactionId:   req.TransactionId,
		TransactionType: req.TransactionType,
		FromAccountId:   req.FromAccountId,
		ToAccountId:     req.ToAccountId,
		Amount:          req.Amount,
		Currency:        req.Currency,
		ExpansionFactor: req.ExpansionFactor,
	}
	//开启事务
	transactionErr := mysql_dao.DB.Transaction(func(tx *gorm.DB) error {
		// sql get
		fromErr := mysql_dao.GetBalance(ctx, &fromAccount, req.FromAccountId)
		toErr := mysql_dao.GetBalance(ctx, &toAccount, req.ToAccountId)
		// AccountId not found
		if errors.Is(fromErr, gorm.ErrRecordNotFound) {
			return utils.ErrFromAccountNotFound
		}
		if errors.Is(toErr, gorm.ErrRecordNotFound) {
			return utils.ErrToAccountNotFound
		}
		// SQL 执行出错（非“记录不存在”）
		if fromErr != nil {
			return fromErr
		}
		if toErr != nil {
			return fromErr
		}

		// compar balence
		if fromAccount.Balance < req.Amount {
			return utils.ErrInsufficientBalance
		}
		// calculation
		fromAccount.Balance = fromAccount.Balance - req.Amount
		toAccount.Balance = toAccount.Balance + req.Amount
		fromAccount.Version++
		toAccount.Version++

		// 顺序上锁id
		var (
			small mysql_model.Balance
			big   mysql_model.Balance
		)
		if fromAccount.BalanceAccountId < toAccount.BalanceAccountId {
			small = fromAccount
			big = toAccount
		} else {
			small = toAccount
			big = fromAccount
		}
		log.Println(small.BalanceAccountId, big.BalanceAccountId)

		var mysqlErr *mysql.MySQLError
		// CreateTransfer
		err := mysql_dao.CreateTransfer(ctx, tx, &transfer)
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return utils.ErrDuplicateRequest
			//case 1045:
			//	return utils.ErrAccountNotFound
			// todo test controller也没写，这两种什么场景？
			case 1048:
				return utils.ErrFieldZeroOrNull
			case 1406:
				return utils.ErrFieldLength
			}
		}

		// Update small Balance
		log.Println("UpdateTransfer ing....", time.Now())
		rows, err := mysql_dao.UpdateTransferBalance(ctx, tx, &small)
		if err != nil {
			return err
		}
		if rows == 0 {
			return utils.ErrOptimisticLockConflict
		}
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1205 {
				return utils.ErrLockWaitTimeOut
			}
		}
		// Update big Balance
		rows, err = mysql_dao.UpdateTransferBalance(ctx, tx, &big)
		if err != nil {
			return err
		}
		if rows == 0 {
			return utils.ErrOptimisticLockConflict
		}
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1205 {
				return utils.ErrLockWaitTimeOut
			}
		}

		return err
	})

	// 事务失败
	if transactionErr != nil {
		return nil, transactionErr
	}

	// redis set A
	fromErr := mysql_dao.GetBalance(ctx, &fromAccount, req.FromAccountId)
	toErr := mysql_dao.GetBalance(ctx, &toAccount, req.ToAccountId)
	if fromErr != nil {
		log.Println(fromErr)
	}
	if toErr != nil {
		log.Println(toErr)
	}
	redisFromAccountA := &redis_model.Balance{
		BalanceAccountId: fromAccount.BalanceAccountId,
		Balance:          fromAccount.Balance,
		CreateTime:       fromAccount.CreateTime,
		UpdateTime:       fromAccount.UpdateTime,
		Currency:         fromAccount.Currency,
	}
	err := redis_dao.SetBalanceCache(ctx, redisFromAccountA.BalanceAccountId, redisFromAccountA)
	if err != nil {
		log.Println("redis set error!")
	}
	redisFromAccountB := &redis_model.Balance{
		BalanceAccountId: toAccount.BalanceAccountId,
		Balance:          toAccount.Balance,
		CreateTime:       toAccount.CreateTime,
		UpdateTime:       toAccount.UpdateTime,
		Currency:         toAccount.Currency,
	}
	err = redis_dao.SetBalanceCache(ctx, redisFromAccountB.BalanceAccountId, redisFromAccountB)
	if err != nil {
		log.Println("redis set error!")
	}

	// get Transfer
	err = mysql_dao.GetTransfer(ctx, mysql_dao.DB, &transfer)
	if err != nil {
		log.Println("transfer set error!")
	}

	// return
	return &model.TransferBalanceResp{
		TransactionId:   transfer.TransactionId,
		TransactionType: transfer.TransactionType,
		FromAccountId:   transfer.FromAccountId,
		ToAccountId:     transfer.ToAccountId,
		Amount:          transfer.Amount,
		Currency:        transfer.Currency,
		CreateTime:      transfer.CreateTime,
		UpdateTime:      transfer.UpdateTime,
		ExpansionFactor: transfer.ExpansionFactor,
	}, nil
}
