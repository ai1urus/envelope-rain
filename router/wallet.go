package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.WithFields(log.Fields{
		"uid": uid,
	}).Info("User check wallet")
	// log.Infof("query %s's wallet", uid)

	envelope_list, err := rdb.SMembers(fmt.Sprintf("EnvelopeList:%v", uid)).Result()
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
	amount, _ := rdb.HGet(fmt.Sprintf("UserInfo:%v", uid), "amount").Result()

	envelopes := []gin.H{}

	for _, eid := range envelope_list {
		_envelope, _ := rdb.HGetAll(fmt.Sprintf("EnvelopeInfo:%v", eid)).Result()
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
