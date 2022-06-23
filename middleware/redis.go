package middleware

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client
var cfg config.RedisConfig
var once sync.Once

func loadUserInfo() {
	var users []database.User

	result := database.GetDB().Find(&users)
	if result.Error != nil {
		panic(fmt.Sprintf("Load UserInfo from DB failed! error: %v", result.Error))
	}

	// var err error
	pipe := rdb.Pipeline()
	for _, user := range users {
		pipe.Set(fmt.Sprintf("UserCount:%v", user.Uid), user.Cur_count, time.Duration(10)*time.Minute)
		pipe.Set(fmt.Sprintf("UserValue:%v", user.Uid), user.Amount, time.Duration(20)*time.Minute)
		// rdb.HMSet(fmt.Sprintf("UserInfo:%v", user.Uid), map[string]interface{}{
		// 	"amount":    user.Amount,
		// 	"cur_count": user.Cur_count})
	}
	_, err := pipe.Exec()
	if err != nil {
		panic(fmt.Sprintf("Redis write UserInfo failed! error: %v", err))
	}
	// fmt.Printf("Redis init: %v", ret)
}

func CreateRedisClient() {
	cfg = config.GetRedisConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})
}

func InitRedis() {
	// CreateRedisClient()

	cfg = config.GetRedisConfig()

	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})

	// init envelope id
	err := rdb.Set("LastEnvelopeID", "0", 0).Err()

	if err != nil {
		panic(err)
	}

	loadUserInfo()
}

func GetRedis() *redis.Client {
	// once.Do(func() {
	// 	InitRedis()
	// })

	return rdb
}
