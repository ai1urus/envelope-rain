package middleware

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

// var ctx = context.Background()

func TestRedisBasic(t *testing.T) {
	InitRedis()

	err := rdb.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func TestRedisWarmup(t *testing.T) {
	var users []database.User
	result := database.GetDB().Find(&users)

	if result.Error != nil {
		panic(result.Error)
	}

	for _, user := range users {
		GetRedis().HMSet(fmt.Sprintf("UserInfo:%v", user.Uid), map[string]interface{}{"amount": user.Amount, "cur_count": user.Cur_count})
	}
}

func TestRedisGetHash(t *testing.T) {
	user, err := GetRedis().HMGet("UserInfo:1", "amount", "cur_count").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

func TestRedisAddSet(t *testing.T) {
	// GetRedis().SAdd()
}

func TestRedisGetHashAll(t *testing.T) {
	InitRedis()
	result, err := GetRedis().HGetAll("EnvelopeInfo:0").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestRedisGetSet(t *testing.T) {
	InitRedis()
	result, err := GetRedis().HGetAll("EnvelopeInfo:0").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestRedisIncr(t *testing.T) {
	InitRedis()
	var wg sync.WaitGroup

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				fmt.Println(GetRedis().Incr("LastEid").Val())
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("LastEid is %v\n", GetRedis().Incr("LastEid").Val())
}

func TestFLoatMultiInt(t *testing.T) {
	var ia int = 100
	fmt.Println(float64(ia) * 0.1)
}

func TestIncrBy(t *testing.T) {
	config.InitConfig()
	CreateRedisClient()
	// rdb.Set("UserCount", 0, 1000000)
	result, err := rdb.IncrBy("UserCount", 100).Result()
	fmt.Println(result, err)
}

func TestSingleSet(t *testing.T) {
	config.InitConfig()
	CreateRedisClient()
	// rdb.Set("UserCount", 0, 1000000)
	err := rdb.Set("LastEnvelopeId", "0", 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

func TestReadNotExist(t *testing.T) {
	config.InitConfig()
	CreateRedisClient()
	// rdb.Set("UserCount", 0, 1000000)
	result, err := rdb.Get("UserInfo:999888999").Int()
	if err == redis.Nil {
		fmt.Println("Value Not exist!")
	} else if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestRedisInitForRed(t *testing.T) {
	config.InitConfig()
	CreateRedisClient()

	pipe := rdb.Pipeline()
	for i := 0; i < 100010; i++ {
		pipe.Set(fmt.Sprintf("User:%v:Snatch", i), 0, time.Duration(10)*time.Minute)
	}
	ret, err := pipe.Exec()
	rdb.Set("TotalMoney", 1000000000000, 0)
	rdb.Set("MaxCount", 5, 0)
	rdb.Set("Probability", 100, 0)
	rdb.Set("EnvelopeNum", 100000000, 0)
	fmt.Println(ret, err)
}
