package util

import (
	"envelope-rain/config"
	"envelope-rain/database"
	"envelope-rain/middleware"
	"strconv"
	"time"

	"math/rand"
)

func GetEnvelopeValue(remain_money int64, min_money, max_money, remain_envelope int32) (int64, int64) {
	_eid, _ := middleware.GetRedis().Get("LastEnvelopeId").Result()
	eid, _ := strconv.ParseInt(_eid, 10, 64)

	err2 := middleware.GetRedis().Set("LastEnvelopeId", strconv.FormatInt(eid+1, 10), -1).Err()
	if err2 != nil {
		panic(err2)
	}

	// fmt.Println(eid)

	if remain_envelope == 1 {
		return eid, Min(remain_money, max_money).(int64)
	}
	// 截尾正态分布，以mean_money为均值，截断范围min_money~max_money
	mean_money := int32(remain_money / int64(remain_envelope))
	// fmt.Println(mean_money)
	max_money = Min(max_money, 2*mean_money-min_money).(int32)
	// fmt.Println(max_money)
	// fmt.Println(max_money - min_money + 1)
	money := min_money + rand.Int31n(max_money-min_money+1)

	return eid, int64(money)
}

func GetEnvelope() (envelope database.Envelope) {
	config := config.GetConfig()
	remain_money := config.GetInt64("envelope.total_money")
	remain_envelope := config.GetInt32("envelope.total_envelope")
	max_money := config.GetInt32("envelope.max_money")
	min_money := config.GetInt32("envelope.min_money")

	if remain_money > 0 && remain_envelope > 0 {
		envelope.Eid, envelope.Value = GetEnvelopeValue(remain_money, min_money, max_money, remain_envelope)
		config.Set("envelope.total_money", remain_money-envelope.Value)
		config.Set("envelope.total_envelope", remain_envelope-1)
	} else {
		envelope.Eid, envelope.Value = -1, 0
	}

	envelope.Opened = false
	envelope.Snatch_time = time.Now().Unix()

	return
}
