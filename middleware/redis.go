package middleware

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var rdb *redis.Client
var cfg config.RedisConfig
var once sync.Once

func loadUserInfo() {
	var users []database.User
	result := database.GetDB().Find(&users)

	if result.Error != nil {
		panic(result.Error)
	}

	for _, user := range users {
		rdb.HMSet(fmt.Sprintf("UserInfo:%v", user.Uid), map[string]interface{}{
			"amount":    user.Amount,
			"cur_count": user.Cur_count})
	}
}

func InitRedis() {
	cfg = config.GetRedisConfig()
	_rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})

	// init envelope id
	err := _rdb.Set("LastEnvelopeId", "0", -1).Err()

	if err != nil {
		panic(err)
	}

	rdb = _rdb

	loadUserInfo()
}

func GetRedis() *redis.Client {
	// once.Do(func() {
	// 	InitRedis()
	// })

	return rdb
}
