package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	// log.Infof("envelope %s opened by %s", envelope_id, uid)

	// 1. Check Eid exist
	envelope, err := rdb.HGetAll(fmt.Sprintf("EnvelopeInfo:%v", envelope_id)).Result()
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
	user, err := rdb.HGetAll(fmt.Sprintf("UserInfo:%v", uid)).Result()
	// Update UserInfo
	user["amount"] = user["amount"] + envelope["value"]
	rdb.HSet(fmt.Sprintf("UserInfo:%v", uid), "amount", user["amount"])
	// Update EnvelopeInfo
	rdb.HSet(fmt.Sprintf("EnvelopeInfo:%v", envelope_id), "opened", true)

	log.WithFields(log.Fields{
		"uid": uid,
		"eid": envelope_id,
	}).Info("User opened")

	// _uid, _ := strconv.ParseInt(uid, 10, 64)
	// value, _ := rdb.Get(envelope_id).Result()
	// database.GetDB().Model(&database.User{}).Where("uid = ?", _uid).Update("amount", value)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": envelope["value"],
		},
	})
}
