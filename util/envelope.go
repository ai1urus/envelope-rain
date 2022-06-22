package util

import (
	"envelope-rain/config"
	"envelope-rain/middleware"
	"sync"
	"sync/atomic"
	"time"

	"math/rand"

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
	valueCache       [][]int64             // 共享变量，更新时锁定单个Cache
	valueCacheSize   int                   // 固定变量
	valueCachePos    int32                 // 共享变量，Generate函数修改
	nextReady        bool                  // 共享变量，二值，用于判断Cache是否锁定
	updateRunning    int32                 // 锁，控制单个Cache的更新锁定
	valueCacheRWLock sync.RWMutex          // 锁，控制单个Cache的读写锁定
}

var eg *EnvelopeGenerator

// 避免锁的使用, 生成一个Cache的Value数据
func (eg *EnvelopeGenerator) GenerateEnvelopeValueNoLock() {
	nextCacheId := (eg.valueCacheId + 1) % 2
	var wg sync.WaitGroup
	wg.Add(eg.valueCacheSize)
	for i := 0; i < eg.valueCacheSize; i++ {
		value := GetEnvelopeValue(eg.cfg.TotalMoney, int32(eg.cfg.MinMoney), int32(eg.cfg.MaxMoney), int32(eg.cfg.TotalEnvelope))
		eg.valueCache[nextCacheId][i] = value
		eg.cfg.TotalMoney -= value
		eg.cfg.TotalEnvelope--
		wg.Done()
	}
	wg.Wait()
	eg.nextReady = true
}

// 切换Cache
func (eg *EnvelopeGenerator) SwitchCacheNoLock() {
	// 这两个操作需要在一步里完成
	eg.valueCacheId = (eg.valueCacheId + 1) % 2
	eg.valueCachePos = -1
	//
	eg.nextReady = false
}

func InitEnvelopeGenerator() {
	eg = &EnvelopeGenerator{
		cfg:            config.GetEnvelopeConfig(),
		valueCacheId:   0,
		valueCacheSize: 20000,
		valueCachePos:  -1,
		nextReady:      false,
	}

	eg.valueCache = make([][]int64, 2)
	eg.valueCache[0] = make([]int64, eg.valueCacheSize)
	eg.valueCache[1] = make([]int64, eg.valueCacheSize)

	eg.GenerateEnvelopeValueNoLock()
	eg.SwitchCacheNoLock()
	eg.GenerateEnvelopeValueNoLock()
}

func GetEnvelopeGenerator() *EnvelopeGenerator {
	return eg
}

func GetEnvelopeValue(remain_money int64, min_money, max_money, remain_envelope int32) int64 {
	if remain_envelope == 1 {
		return Min(remain_money, max_money).(int64)
	}
	// 截尾正态分布，以mean_money为均值，截断范围min_money~max_money
	mean_money := int32(remain_money / int64(remain_envelope))
	max_money = Min(max_money, 2*mean_money-min_money).(int32)
	money := min_money + rand.Int31n(max_money-min_money+1)

	return int64(money)
}

func (eg *EnvelopeGenerator) GetEnvelope() (int64, int64) {
	// return middleware.GetRedis().Incr("LastEnvelopeID").Val(), 1
	for true {
		eg.valueCacheRWLock.RLock()

		eid := middleware.GetRedis().Incr("LastEnvelopeID").Val()
		// try
		pos := atomic.AddInt32(&eg.valueCachePos, 1)
		// CAS锁, 更新cache, 更新时机为nextReady为false且当前cache使用超过10%
		if !eg.nextReady && (pos > int32(0.1*float64(eg.valueCacheSize)-1)) && atomic.CompareAndSwapInt32(&eg.updateRunning, 0, 1) {
			go func() {
				eg.GenerateEnvelopeValueNoLock()
				eg.valueCacheRWLock.Lock()
				eg.nextReady = true
				eg.updateRunning = 0
				eg.valueCacheRWLock.Unlock()
			}()
		}

		pos = atomic.AddInt32(&eg.valueCachePos, 1)
		if pos < int32(eg.valueCacheSize) {
			value := eg.valueCache[eg.valueCacheId][pos]
			eg.valueCacheRWLock.RUnlock()
			return eid, value
		}

		eg.valueCacheRWLock.RUnlock()

		time.Sleep(10)

		eg.valueCacheRWLock.Lock()

		pos = atomic.AddInt32(&eg.valueCachePos, 1)
		if pos < int32(eg.valueCacheSize) {
			value := eg.valueCache[eg.valueCacheId][pos]
			eg.valueCacheRWLock.Unlock()
			return eid, value
		}

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
