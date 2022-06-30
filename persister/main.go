package main

import (
	"context"
	"envelope-rain/config"
	"envelope-rain/database"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Consumer start!")
	config.InitConfig()
	database.InitDB()
	// db := database.GetDB()
	client, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"10.214.150.171:9876"})),
		consumer.WithRetry(2),
		consumer.WithNamespace("ENVELOPE"),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: "rocketmq2",
			SecretKey: "12345678",
		}),
		consumer.WithGroupName("GIT_Group"),
	)
	if err != nil {
		fmt.Println("Init consumer error: " + err.Error())
	}

	err = client.Subscribe("Msg", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := 0; i < len(msgs); i++ {
			fmt.Println(msgs[i])
			switch string(msgs[i].Body) {
			case "CREATE_ENVELOPE":
				// 查db是否存在
				eid, _ := strconv.ParseInt(msgs[i].GetKeys(), 10, 64)
				_, err := database.GetEnvelopeByEid(eid)
				prop := msgs[i].GetProperties()

				if errors.Is(err, gorm.ErrRecordNotFound) {
					var uid, value, snatch_time int64

					uid, err = strconv.ParseInt(prop["uid"], 10, 64)
					value, err = strconv.ParseInt(prop["value"], 10, 64)
					snatch_time, err = strconv.ParseInt(prop["snatch_time"], 10, 64)

					envelope := database.Envelope{
						Eid:         eid,
						Uid:         uid,
						Value:       value,
						Opened:      false,
						Snatch_time: snatch_time,
					}

					database.CreateEnvelope(envelope)

					database.UpdateUserCount(uid)
				}
			case "OPEN_ENVELOPE":
				eid, _ := strconv.ParseInt(msgs[i].GetKeys(), 10, 64)
				envelope, err := database.GetEnvelopeByEid(eid)
				// prop := msgs[i].GetProperties()

				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 重试消息
					panic("Persister write error!")
				} else {
					// 开启
					database.SetEnvelopeOpen(envelope.Eid)

					database.UpdateUserValue(envelope.Uid, envelope.Value)
				}
			}
		}

		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Hour)
	err = client.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}
