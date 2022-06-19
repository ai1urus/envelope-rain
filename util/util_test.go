package util

import (
	"envelope-rain/middleware"
	"fmt"
	"sync"
	"testing"
)

func TestGetEnvelope(t *testing.T) {
	middleware.InitRedis()
	var wg sync.WaitGroup
	wg.Add(5)
	InitEnvelopeGenerator()

	for i := 0; i < 5; i++ {
		go func() {
			sum := 0
			for j := 0; j < 1000; j++ {
				envelope := GetEnvelopeGenerator().GetEnvelope()
				// fmt.Println(envelope)
				sum += int(envelope.Value)
			}
			fmt.Println(sum)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Printf("left envelope is %v", eg.config.remainEnvelope)
	fmt.Printf("left money is %v", eg.config.remainMoney)
	fmt.Println("All done")
}
