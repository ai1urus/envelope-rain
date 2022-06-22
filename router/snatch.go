package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp: true,
	// })

	// 1. 【本地缓存】布隆过滤器判断是否已经达到MaxCount

	// 2. 【Redis】获取CurCount
	// rdb.IncrBy()
	// rdb.IncrBy()
	cur_count, err := rdb.HGet("UserInfo:"+uid, "cur_count").Int()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Redis connect failed",
		})
		return
	}

	// 3. 判断是否达到MaxCount，如果达到则加入布隆过滤器
	max_count := cfg.MaxCount
	if cur_count >= max_count {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "user snatch reached limit",
		})
		return
	}

	// 3. 检查剩余红包数量TODO
	eid, value := eg.GetEnvelope()
	snatch_time := time.Now().Unix()

	if eid == -1 {
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// }).Info("No envelope left")

		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "No envelope left",
		})
		return
	}

	// Redis New EnvelopeInfo
	rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", eid), map[string]interface{}{
		"uid":         uid,
		"value":       value,
		"opened":      false,
		"snatch_time": snatch_time,
	})

	// Redis Update UserList
	rdb.SAdd(fmt.Sprintf("EnvelopeList:%v", uid), eid)

	// Redis Update UserInfo
	cur_count++
	rdb.HSet("UserInfo:"+uid, "cur_count", cur_count)

	// log.WithFields(log.Fields{
	// 	"uid": uid,
	// }).Info("User snatched")

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"envelope_id": eid,
			"max_count":   max_count,
			"cur_count":   cur_count,
		},
	})
}
