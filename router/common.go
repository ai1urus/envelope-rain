package router

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"envelope-rain/middleware"
	"envelope-rain/util"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/bits-and-blooms/bloom"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var cfg config.CommonConfig
var eg *util.EnvelopeGenerator
var mqp rocketmq.Producer

// var server RainServer
var server APIServer

type APIServer struct {
	sendall     bool
	bloomFilter *bloom.BloomFilter
}

func InitService() {
	server.sendall = false
	server.bloomFilter = bloom.NewWithEstimates(1000000, 0.001)
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
	// Init Rocketmq
	middleware.InitProducer()
	mqp = middleware.GetProducer()
	// Init EnvelopeGenerator
	fmt.Println("Init Envelope Generator...")
	util.InitEnvelopeGenerator()
	eg = util.GetEnvelopeGenerator()

	fmt.Println("Init Complete.")
}
