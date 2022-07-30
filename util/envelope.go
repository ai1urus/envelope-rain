package util

import (
	"envelope-rain/config"
	"envelope-rain/middleware"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"math/rand"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// type Segment struct {
// 	IDValue int64
// 	IDMax   int64
// 	IDStep  int64
// }

// type SegmentBuffer struct {
// 	segments      []Segment
// 	currentPos    int
// 	nextReady     bool
// 	initOk        bool
// 	threadRunning bool
// 	lock          sync.Mutex
// }

// type IDGenerator struct {
// 	buffer SegmentBuffer
// }

type EnvelopeGenerator struct {
	cfg              config.EnvelopeConfig // 共享变量：TotalMoney，TotalEnvelope，使用atomic
	valueCacheId     int                   // 共享变量，使用？
	valueCache       [][]int32             // 共享变量，更新时锁定单个Cache
	valueCacheSize   int                   // 固定变量
	valueCachePos    int32                 // 共享变量，Generate函数修改
	nextReady        bool                  // 共享变量，二值，用于判断Cache是否锁定
	updateRunning    int32                 // 锁，控制单个Cache的更新锁定
	valueCacheRWLock sync.RWMutex          // 锁，控制单个Cache的读写锁定
	rdb              *redis.Client
}

var eg *EnvelopeGenerator

func (eg *EnvelopeGenerator) RequestEnvelopeNextBatch() {
	// 申请分布式锁
	for !eg.rdb.SetNX("Lock", 1, 5*time.Minute).Val() {
		time.Sleep(10 * time.Millisecond)
	}

	total_money := eg.rdb.Get("TotalMoney").Val()
	total_envelope := eg.rdb.Get("TotalEnvelope").Val()

	total_money_int, _ := strconv.ParseInt(total_money, 10, 64)
	total_envelope_int, _ := strconv.ParseInt(total_envelope, 10, 64)

	if total_money_int < 1000000 {
		eg.cfg.TotalMoney = total_money_int
	} else {
		eg.cfg.TotalMoney = 1000000
	}

	if total_envelope_int < 5000 {
		eg.cfg.TotalEnvelope = total_envelope_int
	} else {
		eg.cfg.TotalEnvelope = 5000
	}

	total_money_int -= eg.cfg.TotalMoney
	total_envelope_int -= eg.cfg.TotalEnvelope

	eg.rdb.Set("TotalMoney", total_money_int, 0)
	eg.rdb.Set("TotalEnvelope", total_envelope_int, 0)

	eg.rdb.Expire("Lock", 0)
}

// 避免锁的使用, 生成一个Cache的Value数据
func (eg *EnvelopeGenerator) GenerateEnvelopeValueNoLock() {

	nextCacheId := (eg.valueCacheId + 1) % 2

	flag := false

	for i := 0; i < eg.valueCacheSize; i++ {
		if !flag && (eg.cfg.TotalMoney == 0 || eg.cfg.TotalEnvelope == 0) {
			eg.RequestEnvelopeNextBatch()
			if eg.cfg.TotalMoney == 0 || eg.cfg.TotalEnvelope == 0 {
				flag = true
			}
		}

		value := GetEnvelopeValue(eg.cfg.TotalMoney, int32(eg.cfg.MinMoney), int32(eg.cfg.MaxMoney), int32(eg.cfg.TotalEnvelope))
		eg.valueCache[nextCacheId][i] = value
		eg.cfg.TotalMoney -= int64(value)
		eg.cfg.TotalEnvelope--
	}

	eg.nextReady = true
}

// 切换Cache
func (eg *EnvelopeGenerator) SwitchCacheNoLock() {
	// fmt.Println("Change value cache")
	eg.valueCacheId = (eg.valueCacheId + 1) % 2
	eg.valueCachePos = -1
	eg.nextReady = false
}

func InitEnvelopeGenerator() {
	eg = &EnvelopeGenerator{
		cfg:            config.GetEnvelopeConfig(),
		valueCacheId:   0,
		valueCacheSize: 20000,
		valueCachePos:  -1,
		nextReady:      false,
		rdb:            middleware.GetRedis(),
	}

	// 测试分布式锁
	eg.cfg.TotalMoney = 0
	eg.cfg.TotalEnvelope = 0

	eg.valueCache = make([][]int32, 2)
	eg.valueCache[0] = make([]int32, eg.valueCacheSize)
	eg.valueCache[1] = make([]int32, eg.valueCacheSize)

	eg.GenerateEnvelopeValueNoLock()
	// eg.GetNotUsedMoneyJustForTest()
	eg.SwitchCacheNoLock()
	eg.GenerateEnvelopeValueNoLock()
	// eg.GetNotUsedMoneyJustForTest()
}

func GetEnvelopeGenerator() *EnvelopeGenerator {
	return eg
}

func GetEnvelopeValue(remain_money int64, min_money, max_money, remain_envelope int32) int32 {
	if remain_envelope <= 0 || remain_money <= 0 {
		return 0
	}
	if remain_envelope == 1 {
		return Min(int32(remain_money), max_money)
	}
	// 截尾正态分布?以mean_money为均值，截断范围min_money~max_money
	mean_money := int32(remain_money / int64(remain_envelope))
	max_money = Min(max_money, 2*mean_money-min_money)
	money := min_money + rand.Int31n(max_money-min_money+1)

	return money
}

func (eg *EnvelopeGenerator) GetEnvelope() (int64, int32) {
	// // return middleware.GetRedis().Incr("LastEnvelopeID").Val(), 1
	// eid := eg.rdb.Incr("LastEnvelopeID").Val()
	// // if eg.cfg.TotalMoney == 0 || eg.cfg.TotalEnvelope == 0 {
	// eg.RequestEnvelopeNextBatch()
	// // }
	// value := GetEnvelopeValue(100000, int32(eg.cfg.MinMoney), int32(eg.cfg.MaxMoney), 1000)
	// return eid, value

	for true {
		eg.valueCacheRWLock.RLock()
		// try
		// pos := atomic.AddInt32(&eg.valueCachePos, 1)
		// CAS锁, 更新cache, 更新时机为nextReady为false且当前cache使用超过10%
		if !eg.nextReady && (atomic.LoadInt32(&eg.valueCachePos) > int32(0.1*float64(eg.valueCacheSize)-1)) && atomic.CompareAndSwapInt32(&eg.updateRunning, 0, 1) {
			go func() {
				eg.GenerateEnvelopeValueNoLock()
				eg.valueCacheRWLock.Lock()
				eg.nextReady = true
				eg.updateRunning = 0
				eg.valueCacheRWLock.Unlock()
			}()
		}

		pos := atomic.AddInt32(&eg.valueCachePos, 1)
		if pos < int32(eg.valueCacheSize) {
			value := eg.valueCache[eg.valueCacheId][pos]
			eg.valueCacheRWLock.RUnlock()
			eid := eg.rdb.Incr("LastEnvelopeID").Val()
			return eid, value
		}

		eg.valueCacheRWLock.RUnlock()

		// 当前读到的pos超出限度，需要进行buffer切换
		time.Sleep(10)

		eg.valueCacheRWLock.Lock()

		// 后续的pos超出范围的goroutine先获取一个新pos，以此判断buffer是否切换完成
		pos = atomic.AddInt32(&eg.valueCachePos, 1)
		if pos < int32(eg.valueCacheSize) {
			value := eg.valueCache[eg.valueCacheId][pos]
			eg.valueCacheRWLock.Unlock()
			eid := eg.rdb.Incr("LastEnvelopeID").Val()
			return eid, value
		}

		// 第一个拿到写锁的goroutine执行buffer切换
		if eg.nextReady {
			eg.SwitchCacheNoLock()
		} else {
			log.Error("Both two envelope buffer empty!")
			eg.valueCacheRWLock.Unlock()
			break
		}

		eg.valueCacheRWLock.Unlock()
	}

	return -1, -1
}

func (eg *EnvelopeGenerator) JustForTestGetUsedEnvelope() (count, value int64) {
	if eg.valueCachePos < int32(eg.valueCacheSize) {
		fmt.Println("")
		countA := int64(eg.valueCacheSize) - int64(eg.valueCachePos) - 1
		var valueA int32 = 0
		for i := eg.valueCachePos + 1; i < int32(eg.valueCacheSize); i++ {
			valueA += eg.valueCache[eg.valueCacheId][i]
		}
		fmt.Printf("当前Buffer Count: %v, Value: %v\n", countA, valueA)
		count += countA
		value += int64(valueA)
	}
	// return count, value
	if eg.nextReady {
		var tmpId int = (eg.valueCacheId + 1) % 2
		countB := int64(eg.valueCacheSize)
		var valueB int32 = 0
		for i := 0; i < eg.valueCacheSize; i++ {
			valueB += eg.valueCache[tmpId][i]
		}
		fmt.Printf("缓冲Buffer Count: %v, Value: %v\n", countB, valueB)
		count += countB
		value += int64(valueB)
	}

	return count, value
}
