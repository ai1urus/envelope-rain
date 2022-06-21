package middleware

import (
	"envelope-rain/database"
	"fmt"
	"sync"
	"testing"

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
