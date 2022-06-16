package main

import (
	"envelope-rain/util"

	"fmt"
	"testing"
	"time"
)

func TestTimeStamp(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestMin(t *testing.T) {
	var a int64 = 10000
	var b int32 = 1999

	fmt.Println(util.Min(a, b))
}
