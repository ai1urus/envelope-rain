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

	// logic start
	// 1. 检查用户是否存在
	cur_count, err := rdb.HGet(fmt.Sprintf("UserInfo:%v", uid), "cur_count").Int()
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// }).Info("User not found")

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "user not found",
		})
		return
	}

	// 2. 检查已抢次数是否超出
	// fmt.Println(cfg)
	max_count := cfg.MaxCount
	// fmt.Printf("cur_count %v max_count %v\n", cur_count, max_count)
	if cur_count >= max_count {
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// }).Info("User snatch reached limit")

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
	rdb.HSet(fmt.Sprintf("UserInfo:%v", uid), "cur_count", cur_count)

	// log.WithFields(log.Fields{
	// 	"uid": uid,
	// }).Info("User snatched")
	// fmt.Println("User snatched")

	// return

	// logic end

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
