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
	message := primitive.NewMessage(topic, []byte("CREATE_ENVELOPE"))
	message.WithKeys([]string{"100"})
	fmt.Println(message.GetKeys())
	message.WithProperties(params)
	wg.Add(1)
	result, err := p.SendSync(context.Background(), message)
	wg.Done()
	// err = p.SendAsync(context.Background(),
	// 	func(ctx context.Context, result *primitive.SendResult, e error) {
	// 		if e != nil {
	// 			fmt.Printf("receive message error: %s\n", err)
	// 		} else {
	// 			fmt.Printf("send message success: result=%s\n", result.String())
	// 		}
	// 		wg.Done()
	// 	}, message)
	fmt.Println(result)
	if err != nil {
		fmt.Printf("SnatchHandler label 9, an error occurred when sending message:%s\n", err)
	}
	wg.Wait()
}

func TestMsgKeys(t *testing.T) {
	InitProducer()
	var err error
	p := GetProducer()
	var wg sync.WaitGroup
	message := primitive.NewMessage("Msg", []byte("CREATE_ENVELOPE"))
	message.WithProperties(map[string]string{
		"eid":         "100",
		"uid":         "100",
		"value":       "100",
		"opened":      "false",
		"snatch_time": "100",
	})
	message.WithKeys([]string{"100"})
	fmt.Println(message.GetKeys())
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

func TestStr2Int(t *testing.T) {
	a := "100"
	eid, err := strconv.ParseInt(a, 10, 64)
	fmt.Println(eid, err)
}
