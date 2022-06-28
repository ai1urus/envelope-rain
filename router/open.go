package router

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	// log.Infof("envelope %s opened by %s", envelope_id, uid)

	// 1. 调用Lua脚本
	result, err := rdb.EvalSha(openHash, []string{uid, envelope_id}).Int()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 11,
			"msg":  "Service inavailabel",
		})
		return
	}

	switch result {
	case -1:
		log.WithFields(log.Fields{
			"eid": envelope_id,
		}).Info("Envelope not exist")

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Envelope not exist",
		})
		return
	case -2:
		log.WithFields(log.Fields{
			"uid": uid,
			"eid": envelope_id,
		}).Info("Uid not match with Eid")

		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "Uid not match with Eid",
		})
		return
	case -3:
		log.WithFields(log.Fields{
			"uid": uid,
			"eid": envelope_id,
		}).Info("Envelope already opened")

		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "Envelope already opened",
		})
		return
	default:
		// 成功打开红包，用户余额增加
		_, err = rdb.IncrBy("UserValue:"+uid, int64(result)).Result()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 11,
				"msg":  "Service inavailabel",
			})
			return
		}
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// 	"eid": envelope_id,
		// }).Info("User opened")

		// _uid, _ := strconv.ParseInt(uid, 10, 64)
		// value, _ := rdb.Get(envelope_id).Result()
		// database.GetDB().Model(&database.User{}).Where("uid = ?", _uid).Update("amount", value)

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"value": result,
			},
		})
	}
}
