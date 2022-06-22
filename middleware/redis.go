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

	// var err error

	for _, user := range users {
		// err = rdb.Set(fmt.Sprintf("UserCount:%v", user.Uid), user.Cur_count, time.Duration(10)*time.Minute).Err()
		// if err != nil {
		// 	panic(err)
		// }
		// err = rdb.Set(fmt.Sprintf("UserBalance:%v", user.Uid), user.Amount, time.Duration(20)*time.Minute).Err()
		// if err != nil {
		// 	panic(err)
		// }
		rdb.HMSet(fmt.Sprintf("UserInfo:%v", user.Uid), map[string]interface{}{
			"amount":    user.Amount,
			"cur_count": user.Cur_count})
	}
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
	err := rdb.Set("LastEnvelopeId", "0", 0).Err()

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
