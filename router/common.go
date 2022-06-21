package router

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"envelope-rain/middleware"
	"envelope-rain/util"
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client
var cfg config.CommonConfig

func InitService() {
	fmt.Println("Init Config...")
	config.InitConfig()
	cfg = config.GetCommonConfig()
	// Init DB
	fmt.Println("Init DB...")
	database.InitDB()
	// Init Redis
	fmt.Println("Init Redis...")
	middleware.InitRedis()
	rdb = middleware.GetRedis()
	// Init EnvelopeGenerator
	fmt.Println("Init Envelope Generator...")
	util.InitEnvelopeGenerator()

	fmt.Println("Init Complete.")
}
