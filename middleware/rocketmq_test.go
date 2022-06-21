package middleware

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func TestRocketmqProduce(t *testing.T) {
	InitProducer()
	var err error
	p := GetProducer()
	var wg sync.WaitGroup
	topic := "Msg"
	params := make(map[string]string)
	params["UID"] = strconv.FormatInt(1234, 10)
	params["EID"] = strconv.FormatInt(1123, 10)
	params["Value"] = strconv.Itoa(12415)
	params["SnatchTime"] = strconv.Itoa(int(time.Now().Unix()))
	message := primitive.NewMessage(topic, []byte("create_envelope"))
	message.WithProperties(params)
	wg.Add(1)
	err = p.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive message error: %s\n", err)
			} else {
				fmt.Printf("send message success: result=%s\n", result.String())
			}
			wg.Done()
		}, message)
	if err != nil {
		fmt.Printf("SnatchHandler label 9, an error occurred when sending message:%s\n", err)
	}
	wg.Wait()
}
