package utils

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrFromAccountNotFound = errors.New("the fromAccount does not exist")
	ErrToAccountNotFound   = errors.New("the fromAccount does not exist")

	ErrDuplicateRequest       = errors.New("this is a repetitive operation")
	ErrAccountFound           = errors.New("the account existed")
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict: data has been updated by another transaction")
	ErrFieldZeroOrNull        = errors.New("the field is 0 or empty")
	ErrFieldLength            = errors.New("the data exceeds the field length")
	ErrLockWaitTimeOut        = errors.New("lock wait timeout")
	//ErrDatabaseBusy           = errors.New("database is busy, please try again later")
	//ErrRelatedRecordNotFound  = errors.New("related record not found") //查找订单用
	//ErrDeadlock               = errors.New("transaction deadlock, please retry")
	//ErrLockNotWait            = errors.New("lock could not be acquired immediately") //悲观锁
)
