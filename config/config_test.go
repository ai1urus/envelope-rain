package config

import (
	"fmt"
	"testing"
)

func TestReadYAML(t *testing.T) {
	InitConfig()
	fmt.Printf("Config test %v\n", cfg.Get("envelope.max_money"))
}
