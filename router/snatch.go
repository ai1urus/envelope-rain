package router

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp: true,
	// })

	// 1. 【本地缓存】布隆过滤器判断是否已经达到MaxCount
	if server.sendall {
		c.JSON(200, gin.H{
			"code": 31,
			"msg":  "No envelope left",
		})
		return
	}

	if server.bloomFilter.TestString(uid) {
		c.JSON(200, gin.H{
			"code": 22,
			"msg":  "User snatch reached limit",
		})
		return
	}

	// 2. 【Redis】获取CurCount
	cur_count, err := rdb.Get("UserCount:" + uid).Int64()
	if err == redis.Nil {
		// Redis运行正常情况下值一定在Redis中，否则说明用户不存在
		c.JSON(200, gin.H{
			"code": 21,
			"msg":  "User not exist",
		})
		return
	} else if err != nil {
		// Redis 连接失败, (可以添加MySQL逻辑继续提供服务?)
		c.JSON(500, gin.H{
			"code": 11,
			"msg":  "Service inavailabel",
		})
		return
	}

	// 判断红包概率
	if rand.Int()%100 > cfg.Probability {
		c.JSON(200, gin.H{
			"code": 23,
			"msg":  "User not lucky",
		})
		return
	}

	// 3. 判断是否达到MaxCount，如果达到则加入布隆过滤器
	cur_count, err = rdb.Incr("UserCount:" + uid).Result()
	max_count := cfg.MaxCount
	if cur_count > max_count {
		server.bloomFilter.AddString(uid)
		c.JSON(200, gin.H{
			"code": 22,
			"msg":  "User snatch reached limit",
		})
		return
	}

	// 4. 抢红包成功
	eid, value := eg.GetEnvelope()
	snatch_time := time.Now().Unix()

	// 红包发完，需要提到前面，localcache加一个flag
	if eid == -1 {
		// log.WithFields(log.Fields{
		// 	"uid": uid,
		// }).Info("No envelope left")
		server.sendall = true
		c.JSON(200, gin.H{
			"code": 31,
			"msg":  "No envelope left",
		})
		return
	}

	_eid := strconv.FormatInt(eid, 10)
	msg := primitive.NewMessage("Msg", []byte("CREATE_ENVELOPE"))
	// msg.WithShardingKey(_eid)
	msg.WithKeys([]string{_eid})
	msg.WithProperties(map[string]string{
		"eid":         _eid,
		"uid":         uid,
		"value":       strconv.Itoa(int(value)),
		"opened":      strconv.FormatBool(false),
		"snatch_time": strconv.FormatInt(snatch_time, 10),
	})

	var wg sync.WaitGroup
	wg.Add(1)
	err = mqp.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				panic(fmt.Sprintf("receive message error: %s\n", err))
			} else {
				// fmt.Printf("send message success: result=%s\n", result.String())
			}
			wg.Done()
		}, msg)

	if err != nil {
		panic(fmt.Sprintf("send message error: %s\n", err))
	}
	wg.Wait()

	// 5. Redis 插入红包
	// rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", eid), envelope)
	rdb.HMSet(fmt.Sprintf("EnvelopeInfo:%v", eid), map[string]interface{}{
		"uid":         uid,
		"value":       value,
		"opened":      false,
		"snatch_time": snatch_time,
	})
	rdb.Expire(fmt.Sprintf("EnvelopeInfo:%v", eid), time.Duration(20)*time.Minute)

	// Redis Update UserList
	rdb.SAdd(fmt.Sprintf("EnvelopeList:%v", uid), eid)

	// Redis Update UserInfo
	// cur_count++
	// rdb.HSet("UserInfo:"+uid, "cur_count", cur_count)

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
