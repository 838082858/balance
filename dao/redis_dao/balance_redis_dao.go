package redis_dao

import (
	"context"
	"encoding/json"
	"http-demo/model/redis_model"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func RedisStorage() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 密码，如果没有则留空
		DB:       0,                // 使用默认的数据库 0,redis有0~15的数据库，之间的数据是隔离的。
	})
	pong, pingErr := rdb.Ping(ctx).Result()
	if pingErr != nil {
		panic(pingErr)
	}
	log.Println("redis connection success:", pong) //成功连接后，pong的值为PONG
}

// GetBalanceCache strconv转换字符串，Format格式，Uint指uint64
func GetBalanceCache(ctx context.Context, balanceAccountId uint64) (*redis_model.Balance, error) {
	val, err := rdb.Get(ctx, strconv.FormatUint(balanceAccountId, 10)).Result()
	var getCacheResult *redis_model.Balance
	unmErr := json.Unmarshal([]byte(val), getCacheResult)
	if unmErr != nil {
		return nil, unmErr
		//一定要在这里转换，行业规定
	}
	return getCacheResult, err
}

// SetBalanceCache 创建新数据
func SetBalanceCache(ctx context.Context, balanceAccountId uint64, value *redis_model.Balance) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, strconv.FormatUint(balanceAccountId, 10), val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteBalanceCache 删除数据
func DeleteBalanceCache(ctx context.Context, balanceAccountId uint64) error {
	err := rdb.Del(ctx, strconv.FormatUint(balanceAccountId, 10)).Err()
	if err != nil {
		return err
	}
	return nil
}
