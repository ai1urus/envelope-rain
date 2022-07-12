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
var openHash string

type APIServer struct {
	sendall     bool
	bloomFilter *bloom.BloomFilter
}

func GenerateOpenScript(rdb *redis.Client) string {
	var openScript string = `
	local uid = KEYS[1]
	local eid = KEYS[2]
	
	local envelope = redis.call("HMGET", "EnvelopeInfo:" .. eid, "uid", "opened", "value")
	
	-- Ret 1 eid 不存在
	if not envelope[1] then
		return -1
	end
	
	-- Ret 2 eid 与 uid 不匹配
	if envelope[1] ~= uid then 
		return -2
	end
	
	-- Ret 3 envelope 已开启
	if envelope[2] == "1" then
		return -3
	end
	
	-- Ret 0 成功打开
	redis.call("HMSET", "EnvelopeInfo:"..eid, "opened", "1") 
	redis.call("INCRBY", "UserValue:"..uid, envelope[3])
	return envelope[3]
	`

	_openHash, err := rdb.ScriptLoad(openScript).Result()
	if err != nil {
		panic(fmt.Sprintf("Open Script create failed: %v", err))
	}
	return _openHash
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
	openHash = GenerateOpenScript(rdb)
	// Init Rocketmq
	// middleware.InitProducer()
	// mqp = middleware.GetProducer()
	// Init EnvelopeGenerator
	fmt.Println("Init Envelope Generator...")
	util.InitEnvelopeGenerator()
	eg = util.GetEnvelopeGenerator()

	fmt.Println("Init Complete.")
}

func StopService() {
	rdb.Close()
}
