package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var globalConfig *viper.Viper
var once sync.Once

func initConfig() {
	fmt.Println("Init GlobalConfig ...")
	globalConfig = viper.New()
	globalConfig.SetConfigName("config")                                     // name of config file (without extension)
	globalConfig.SetConfigType("yaml")                                       // REQUIRED if the config file does not have the extension in the name
	globalConfig.AddConfigPath("/home/ubuntu/Project/envelope-rain/config/") // path to look for the config file in
	globalConfig.AddConfigPath("./config/")
	globalConfig.AddConfigPath(".")    // optionally look for config in the working directory
	err := globalConfig.ReadInConfig() // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

}

func GetConfig() *viper.Viper {
	once.Do(func() {
		initConfig()
	})

	return globalConfig
}
