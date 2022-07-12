package router

import (
	"envelope-rain/database"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	// 1. 调用Lua脚本
	result, err := rdb.EvalSha(openHash, []string{uid, envelope_id}).Int()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "Service unavailable",
		})
		return
	}

	var envelope database.Envelope
	if result == -1 {
		eid, _ := strconv.ParseInt(envelope_id, 10, 64)
		envelope, err = database.GetEnvelopeByEid(eid)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithFields(log.Fields{
				"eid": envelope_id,
			}).Info("Envelope not exist")

			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "Envelope not exist",
			})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{
				"code": 2,
				"msg":  "Service unavailable",
			})
			return
		}

		rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", eid), map[string]interface{}{
			"uid":         envelope.Uid,
			"value":       envelope.Value,
			"opened":      envelope.Opened,
			"snatch_time": envelope.Snatch_time,
		})
		rdb.Expire(fmt.Sprintf("EnvelopeInfo:%v", eid), time.Duration(20)*time.Minute)

		result, err = rdb.EvalSha(openHash, []string{uid, envelope_id}).Int()
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "Service unavailable",
			})
			return
		}
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
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// 	"eid": envelope_id,
		// }).Info("Envelope already opened")

		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "Envelope already opened",
		})
		return
	default:
		// 成功打开红包，用户余额增加
		// msg := primitive.NewMessage("Msg", []byte("OPEN_ENVELOPE"))
		// msg.WithKeys([]string{envelope_id})
		// msg.WithProperties(map[string]string{
		// 	"eid": envelope_id,
		// 	"uid": uid,
		// })

		// _, err = mqp.SendSync(context.Background(), msg)
		// if err != nil {
		// 	// 写 MQ 失败，Count回退？
		// 	c.JSON(500, gin.H{
		// 		"code": 11,
		// 		"msg":  "Service inavailable",
		// 	})
		// 	return
		// }
		// var wg sync.WaitGroup
		// wg.Add(1)
		// err = mqp.SendAsync(context.Background(),
		// 	func(ctx context.Context, result *primitive.SendResult, e error) {
		// 		if e != nil {
		// 			panic(fmt.Sprintf("receive message error: %s\n", err))
		// 		} else {
		// 			// fmt.Printf("send message success: result=%s\n", result.String())
		// 		}
		// 		wg.Done()
		// 	}, msg)
		// if err != nil {
		// 	panic(fmt.Sprintf("send message error: %s\n", err))
		// }
		// wg.Wait()

		// _, err = rdb.IncrBy("UserValue:"+uid, int64(result)).Result()
		// if err != nil {
		// 	c.JSON(500, gin.H{
		// 		"code": 11,
		// 		"msg":  "Service unavailable",
		// 	})
		// 	return
		// }
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
