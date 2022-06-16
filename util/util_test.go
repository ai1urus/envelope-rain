package util

import (
	"envelope-rain/middleware"
	"fmt"
	"testing"
)

func TestGetEnvelope(t *testing.T) {
	middleware.InitRedis()
	fmt.Println(GetEnvelope())
	fmt.Println(GetEnvelope())
	fmt.Println(GetEnvelope())
	fmt.Println(GetEnvelope())
	fmt.Println(GetEnvelope())
}
