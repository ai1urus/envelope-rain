package util

import (
	"math/rand"
	"testing"
	"time"
)

// func TestGetEnvelope(t *testing.T) {
// 	middleware.InitRedis()
// 	var wg sync.WaitGroup
// 	wg.Add(5)
// 	InitEnvelopeGenerator()

// 	for i := 0; i < 5; i++ {
// 		go func() {
// 			sum := 0
// 			for j := 0; j < 1000; j++ {
// 				_, value := GetEnvelopeGenerator().GetEnvelope()
// 				sum += int(value)
// 			}
// 			fmt.Println(sum)
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()

// 	fmt.Printf("left envelope is %v", eg.cfg.TotalEnvelope)
// 	fmt.Printf("left money is %v", eg.cfg.TotalMoney)
// 	fmt.Println("All done")
// }

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
