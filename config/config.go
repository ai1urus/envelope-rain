package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var cfg *viper.Viper
var once sync.Once

type CommonConfig struct {
	TotalMoney    int64
	TotalEnvelope int64
	MaxCount      int64
	MaxMoney      int
	MinMoney      int
	Probability   int
}

type DBConfig struct {
	DBAddr     string
	DBUsername string
	DBPassword string
	DBName     string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type EnvelopeConfig struct {
	TotalMoney    int64
	TotalEnvelope int64
	MaxCount      int
	MaxMoney      int
	MinMoney      int
}

func InitConfig() {
	cfg = viper.New()
	cfg.SetConfigName("config") // name of config file (without extension)
	cfg.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	// cfg.AddConfigPath("/home/jian/Project/envelope-rain/config/") // path to look for the config file in
	cfg.AddConfigPath("./config/")
	cfg.AddConfigPath("../config") // optionally look for config in the working directory
	cfg.AddConfigPath(".")         // optionally look for config in the working directory
	err := cfg.ReadInConfig()      // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func GetCommonConfig() CommonConfig {
	_cfg := CommonConfig{
		TotalMoney:    cfg.GetInt64("envelope.total_money"),
		TotalEnvelope: cfg.GetInt64("envelope.total_envelope"),
		MaxCount:      cfg.GetInt64("envelope.max_snatch"),
		MaxMoney:      cfg.GetInt("envelope.max_money"),
		MinMoney:      cfg.GetInt("envelope.min_money"),
		Probability:   cfg.GetInt("envelope.probability"),
	}
	return _cfg
}

func GetEnvelopeConfig() EnvelopeConfig {
	_cfg := EnvelopeConfig{
		TotalMoney:    cfg.GetInt64("envelope.total_money"),
		TotalEnvelope: cfg.GetInt64("envelope.total_envelope"),
		MaxCount:      cfg.GetInt("envelope.max_snatch"),
		MaxMoney:      cfg.GetInt("envelope.max_money"),
		MinMoney:      cfg.GetInt("envelope.min_money"),
	}
	return _cfg
}

func GetDBConfig() DBConfig {
	_cfg := DBConfig{
		DBAddr:     cfg.GetString("mysql.address"),
		DBUsername: cfg.GetString("mysql.username"),
		DBPassword: cfg.GetString("mysql.password"),
		DBName:     cfg.GetString("mysql.dbname"),
	}
	return _cfg
}

func GetRedisConfig() RedisConfig {
	_cfg := RedisConfig{
		Addr:     cfg.GetString("redis.address"),
		Password: cfg.GetString("redis.password"),
	}
	return _cfg
}
