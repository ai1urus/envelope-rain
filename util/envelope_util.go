package util

import (
	"envelope-rain/config"
	"envelope-rain/middleware"
	"sync"

	"math/rand"
)

type Segment struct {
	IDValue int64
	IDMax   int64
	IDStep  int64
}

type SegmentBuffer struct {
	segments      []Segment
	currentPos    int
	nextReady     bool
	initOk        bool
	threadRunning bool
	lock          sync.Mutex
}

type IDGenerator struct {
	buffer SegmentBuffer
}

type EnvelopeGenerator struct {
	cfg            config.CommonConfig
	currentCache   int
	valueCache     [][]int64
	valueCacheSize int
	valueCachePos  int
	valueCacheLock sync.Mutex
	nextReady      bool
}

var eg *EnvelopeGenerator

func (eg *EnvelopeGenerator) GenerateEnvelopeValue() {
	nextCache := (eg.currentCache + 1) % 2
	for i := 0; i < eg.valueCacheSize; i++ {
		eg.valueCache[nextCache][i] = GetEnvelopeValue(eg.cfg.TotalMoney, int32(eg.cfg.MinMoney), int32(eg.cfg.MaxMoney), int32(eg.cfg.TotalEnvelope))
	}
	eg.nextReady = true
}

func (eg *EnvelopeGenerator) ChangeCache() {
	eg.currentCache = (eg.currentCache + 1) % 2
	eg.valueCachePos = 0
	eg.nextReady = false
}

func InitEnvelopeGenerator() {
	eg = &EnvelopeGenerator{
		cfg:            config.GetCommonConfig(),
		currentCache:   0,
		valueCacheSize: 20000,
		valueCachePos:  0,
		nextReady:      false,
	}

	eg.valueCache = make([][]int64, 2)
	eg.valueCache[0] = make([]int64, eg.valueCacheSize)
	eg.valueCache[1] = make([]int64, eg.valueCacheSize)

	eg.GenerateEnvelopeValue()
	eg.ChangeCache()
	eg.GenerateEnvelopeValue()
	eg.ChangeCache()
}

func GetEnvelopeGenerator() *EnvelopeGenerator {
	return eg
}

func GetEnvelopeValue(remain_money int64, min_money, max_money, remain_envelope int32) int64 {
	// _eid, _ := middleware.GetRedis().Get("LastEnvelopeId").Result()
	// eid, _ := strconv.ParseInt(_eid, 10, 64)
	// err2 := middleware.GetRedis().Set("LastEnvelopeId", strconv.FormatInt(eid+1, 10), -1).Err()
	// if err2 != nil {
	// 	panic(err2)
	// }

	// fmt.Println(eid)

	if remain_envelope == 1 {
		return Min(remain_money, max_money).(int64)
	}
	// 截尾正态分布，以mean_money为均值，截断范围min_money~max_money
	mean_money := int32(remain_money / int64(remain_envelope))
	// fmt.Println(mean_money)
	max_money = Min(max_money, 2*mean_money-min_money).(int32)
	// fmt.Println(max_money)
	// fmt.Println(max_money - min_money + 1)
	money := min_money + rand.Int31n(max_money-min_money+1)

	return int64(money)
}

func (eg *EnvelopeGenerator) GetEnvelope() (int64, int64) {
	eg.valueCacheLock.Lock()
	defer eg.valueCacheLock.Unlock()

	eid := middleware.GetRedis().Incr("LastEnvelopeID").Val()
	value := eg.valueCache[eg.currentCache][eg.valueCachePos]

	eg.valueCachePos++

	if eg.valueCachePos == int(float64(eg.valueCacheSize)*0.1) {
		// 更新buffer
		if !eg.nextReady {
			eg.GenerateEnvelopeValue()
		}
	} else if eg.valueCachePos == eg.valueCacheSize {
		// 切换buffer
		for !eg.nextReady {
			if !eg.nextReady {
				eg.GenerateEnvelopeValue()
			}
		}
		eg.ChangeCache()
	}

	eg.cfg.TotalMoney -= value
	eg.cfg.TotalEnvelope--

	return eid, value
}
