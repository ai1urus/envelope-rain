package config

import (
	"fmt"
	"testing"
)

func TestReadYAML(t *testing.T) {
	conf := GetConfig()
	fmt.Printf("Config test %v\n", conf.Get("envelope.max_money"))
}
