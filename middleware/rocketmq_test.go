package middleware

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func TestRocketmqProduce(t *testing.T) {
	InitProducer()

	p := GetProducer()
	var err error
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		err := p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					fmt.Printf("receive message error: %s\n", err)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, primitive.NewMessage("test", []byte("Hello RocketMQ Go Client!")))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		}
	}
	wg.Wait()
}
