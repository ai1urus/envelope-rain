package main

import (
	"envelope-rain/config"
	"envelope-rain/middleware"
	"envelope-rain/util"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	// init
	// mysql warmup
	// redis warmup
	middleware.InitRedis()
	// router
	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)

	r.Run()
}

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// logic start
	config := config.GetConfig()
	// 1. Check Uid exist
	_cur_count, err := middleware.GetRedis().HGet(fmt.Sprintf("UserInfo:%v", uid), "cur_count").Result()
	if err != nil || _cur_count == "" {
		log.WithFields(log.Fields{
			"uid": uid,
		}).Info("User not found")

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "user not found",
		})
		return
	}

	// 2. Check Count < Limit
	cur_count, _ := strconv.Atoi(_cur_count)
	max_count := config.GetInt("envelope.max_snatch")
	if cur_count >= max_count {
		log.WithFields(log.Fields{
			"uid": uid,
		}).Error("User snatch reached limit")

		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "user snatch reached limit",
		})
		return
	}

	// 3. Chekc remain Envelope
	envelope := util.GetEnvelope()
	if envelope.Eid == -1 {
		log.WithFields(log.Fields{
			"uid": uid,
		}).Info("No envelope left")

		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "No envelope left",
		})
		return
	}

	// Redis New EnvelopeInfo
	middleware.GetRedis().HMSet(fmt.Sprintf("EnvelopeInfo:%v", envelope.Eid), map[string]interface{}{
		"uid":         uid,
		"value":       envelope.Value,
		"opened":      false,
		"snatch_time": envelope.Snatch_time,
	})

	// Redis Update UserList
	middleware.GetRedis().SAdd(fmt.Sprintf("EnvelopeList:%v", uid), envelope.Eid)

	// Redis Update UserInfo
	cur_count++
	middleware.GetRedis().HSet(fmt.Sprintf("UserInfo:%v", uid), "cur_count", cur_count)

	log.WithFields(log.Fields{
		"uid": uid,
	}).Info("User snatched")

	// logic end

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"envelope_id": envelope.Eid,
			"max_count":   max_count,
			"cur_count":   cur_count,
		},
	})
}

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	// log.Infof("envelope %s opened by %s", envelope_id, uid)

	// 1. Check Eid exist
	envelope, err := middleware.GetRedis().HGetAll(fmt.Sprintf("EnvelopeInfo:%v", envelope_id)).Result()
	if err != nil || envelope == nil {
		log.WithFields(log.Fields{
			"eid": envelope_id,
		}).Info("Envelope not exist")

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Envelope not exist",
		})
		return
	}
	// 2. Check Uid match
	if envelope["uid"] != uid {
		log.WithFields(log.Fields{
			"uid": uid,
			"eid": envelope_id,
		}).Info("Uid not match with Eid")

		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "Uid not match with Eid",
		})
		return
	}
	// 3. Check opened
	if envelope["opened"] == "1" {
		log.WithFields(log.Fields{
			"uid": uid,
			"eid": envelope_id,
		}).Info("Envelope already opened")

		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "Envelope already opened",
		})
		return
	}
	// 4. Success open
	user, err := middleware.GetRedis().HGetAll(fmt.Sprintf("UserInfo:%v", uid)).Result()
	// Update UserInfo
	user["amount"] = user["amount"] + envelope["value"]
	middleware.GetRedis().HSet(fmt.Sprintf("UserInfo:%v", uid), "amount", user["amount"])
	// Update EnvelopeInfo
	middleware.GetRedis().HSet(fmt.Sprintf("EnvelopeInfo:%v", envelope_id), "opened", true)

	log.WithFields(log.Fields{
		"uid": uid,
		"eid": envelope_id,
	}).Info("User opened")

	// _uid, _ := strconv.ParseInt(uid, 10, 64)
	// value, _ := middleware.GetRedis().Get(envelope_id).Result()
	// database.GetDB().Model(&database.User{}).Where("uid = ?", _uid).Update("amount", value)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": envelope["value"],
		},
	})
}

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.WithFields(log.Fields{
		"uid": uid,
	}).Info("User check wallet")
	// log.Infof("query %s's wallet", uid)

	envelope_list, err := middleware.GetRedis().SMembers(fmt.Sprintf("EnvelopeList:%v", uid)).Result()
	// 1. Chekc User exist
	if err != nil || envelope_list == nil {
		log.WithFields(log.Fields{
			"uid": uid,
		}).Info("User not found")

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "User not found",
		})
		return
	}
	// 2. Success
	// Redis Get UserInfo
	amount, _ := middleware.GetRedis().HGet(fmt.Sprintf("UserInfo:%v", uid), "amount").Result()

	envelopes := []gin.H{}

	for _, eid := range envelope_list {
		_envelope, _ := middleware.GetRedis().HGetAll(fmt.Sprintf("EnvelopeInfo:%v", eid)).Result()
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

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": envelopes,
		},
	})
}
