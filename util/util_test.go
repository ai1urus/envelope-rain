package util

import (
	"envelope-rain/config"
	"envelope-rain/middleware"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetEnvelope(t *testing.T) {
	config.InitConfig()
	middleware.CreateRedisClient()
	InitEnvelopeGenerator()

	var totalSum int64
	eg := GetEnvelopeGenerator()

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	var sum int64 = 0
	// 	for j := 0; j < 20000; j++ {
	// 		_, value := eg.GetEnvelope()
	// 		sum += int64(value)
	// 	}
	// 	fmt.Println(sum)
	// 	atomic.AddInt64(&totalSum, int64(sum))
	// 	wg.Done()
	// }()
	// wg.Wait()

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			var sum int64 = 0
			for j := 0; j < 45678; j++ {
				_, value := eg.GetEnvelope()
				sum += int64(value)
			}
			atomic.AddInt64(&totalSum, int64(sum))
			wg.Done()
		}()
	}
	wg.Wait()

	count, value := eg.JustForTestGetUsedEnvelope()

	fmt.Printf("used money is %v\n", totalSum+value)
	fmt.Printf("expected used money is %v\n", config.GetCommonConfig().TotalMoney-eg.cfg.TotalMoney)

	usedcount, _ := middleware.GetRedis().Get("LastEnvelopeID").Int64()
	fmt.Printf("used envelope is %v\n", usedcount+count)
	fmt.Printf("expected used envelope is %v\n", config.GetCommonConfig().TotalEnvelope-eg.cfg.TotalEnvelope)
	fmt.Println("All done")
}

func BenchmarkRandom1(b *testing.B) {
	rand.Seed(time.Now().Unix())
	ans := 0
	for i := 0; i < 10000000; i++ {
		ans += rand.Int() % 100
	}
}

func BenchmarkRandom2(b *testing.B) {
	rand.Seed(time.Now().Unix())
	var ans int64 = 0
	for i := 0; i < 10000000; i++ {
		ans += time.Now().UnixNano() % 100
	}
}

func BenchmarkRandom3(b *testing.B) {
	rand.Seed(time.Now().Unix())
	ans := 0
	for i := 0; i < 10000000; i++ {
		ans += rand.Intn(100)
	}
}

func BenchmarkGetEnvelope(b *testing.B) {
	config.InitConfig()
	middleware.CreateRedisClient()
	InitEnvelopeGenerator()
	eg := GetEnvelopeGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eg.GetEnvelope()
	}

	// b.SetParallelism(10)
	// b.RunParallel(func(p *testing.PB) {
	// 	for p.Next() {
	// 		eg.GetEnvelope()
	// 	}
	// })

	// for j := 0; j < 100000; j++ {
	// 	eg.GetEnvelope()
	// }
}
