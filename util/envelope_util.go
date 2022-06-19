package util

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"envelope-rain/middleware"
	"sync"
	"time"

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

type EnvelopeConfig struct {
	remainMoney    int64
	remainEnvelope int32
	maxMoney       int32
	minMoney       int32
}

type EnvelopeGenerator struct {
	config         EnvelopeConfig
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
		eg.valueCache[nextCache][i] = GetEnvelopeValue(eg.config.remainMoney, eg.config.minMoney, eg.config.maxMoney, eg.config.remainEnvelope)
		// config.Set("envelope.total_money", remain_money-envelope.Value)
		// config.Set("envelope.total_envelope", remain_envelope-1)
	}
	eg.nextReady = true
}

func (eg *EnvelopeGenerator) ChangeCache() {
	eg.currentCache = (eg.currentCache + 1) % 2
	eg.valueCachePos = 0
	eg.nextReady = false
}

func InitEnvelopeGenerator() {
	_config := config.GetConfig()

	eg = &EnvelopeGenerator{
		config: EnvelopeConfig{
			remainMoney:    _config.GetInt64("envelope.total_money"),
			remainEnvelope: _config.GetInt32("envelope.total_envelope"),
			maxMoney:       _config.GetInt32("envelope.max_money"),
			minMoney:       _config.GetInt32("envelope.min_money"),
		},
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

func (eg *EnvelopeGenerator) GetEnvelope() (envelope *database.Envelope) {
	eg.valueCacheLock.Lock()
	defer eg.valueCacheLock.Unlock()

	envelope = &database.Envelope{
		Eid:         middleware.GetRedis().Incr("LastEnvelopeID").Val(),
		Value:       eg.valueCache[eg.currentCache][eg.valueCachePos],
		Opened:      false,
		Snatch_time: time.Now().Unix(),
	}

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

	eg.config.remainMoney -= envelope.Value
	eg.config.remainEnvelope--

	// if eg.config.remainMoney > 0 && eg.config.remainEnvelope > 0 {
	// 	envelope.Eid, envelope.Value = GetEnvelopeValue(remain_money, min_money, max_money, remain_envelope)
	// } else {
	// 	envelope.Eid, envelope.Value = -1, 0
	// }

	return
}
