package router

import (
	"envelope-rain/database"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	_uid, _ := strconv.ParseInt(uid, 10, 64)
	// 1. Chekc User exist
	amount, err := rdb.Get("UserValue:" + uid).Int64()
	if err != nil || errors.Is(err, redis.Nil) {
		amount, err = database.GetUserAmount(_uid)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithFields(log.Fields{
				"uid": uid,
			}).Info("User not found")

			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "User not found",
			})
			return
		} else if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "Wallet Service unavailable",
			})
			return
		}
	}

	var envelopes []gin.H
	envelope_list, err := rdb.SMembers("EnvelopeList:" + uid).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		_envelopes, _ := database.GetEnvelopeByUid(_uid)
		if len(_envelopes) > 0 {
			for i := 0; i < len(_envelopes); i++ {
				rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", _envelopes[i].Eid), map[string]interface{}{
					"uid":         uid,
					"value":       _envelopes[i].Value,
					"opened":      _envelopes[i].Opened,
					"snatch_time": _envelopes[i].Snatch_time,
				})
				rdb.Expire(fmt.Sprintf("EnvelopeInfo:%v", _envelopes[i].Eid), time.Duration(20)*time.Minute)
				rdb.SAdd(fmt.Sprintf("EnvelopeList:%v", uid), _envelopes[i].Eid)

				var envelope gin.H = gin.H{
					"envelope_id": _envelopes[i].Eid,
					"opened":      _envelopes[i].Opened,
					"snatch_time": _envelopes[i].Snatch_time,
				}
				if _envelopes[i].Opened == true {
					envelope["value"] = _envelopes[i].Value
				}
				envelopes = append(envelopes, envelope)
			}
		}
	} else {
		for _, eid := range envelope_list {
			_envelope, err := rdb.HGetAll(fmt.Sprintf("EnvelopeInfo:%v", eid)).Result()
			if err != nil {
				_eid, _ := strconv.ParseInt(eid, 10, 64)
				_envelope, err := database.GetEnvelopeByEid(_eid)
				if err != gorm.ErrRecordNotFound {
					rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", _envelope.Eid), map[string]interface{}{
						"uid":         uid,
						"value":       _envelope.Value,
						"opened":      _envelope.Opened,
						"snatch_time": _envelope.Snatch_time,
					})
					rdb.Expire(fmt.Sprintf("EnvelopeInfo:%v", _envelope.Eid), time.Duration(20)*time.Minute)

					var envelope gin.H = gin.H{
						"envelope_id": eid,
						"opened":      _envelope.Opened,
						"snatch_time": _envelope.Snatch_time,
					}
					if _envelope.Opened == true {
						envelope["value"] = _envelope.Value
					}
					envelopes = append(envelopes, envelope)
				}
			} else {
				var envelope gin.H = gin.H{
					"envelope_id": eid,
					"opened":      _envelope["opened"],
					"snatch_time": _envelope["snatch_time"],
				}
				if _envelope["opened"] == "1" {
					envelope["value"] = _envelope["value"]
				}
				envelopes = append(envelopes, envelope)
			}
		}
	}
	// 2. Success

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": envelopes,
		},
	})
}
