package database

import (
	"envelope-rain/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// type Envelope struct {
// 	Envelope_id int64 `gorm:"primaryKey"`
// 	Uid         int64
// 	Value       int64
// 	Opened      bool
// 	Snatch_time int64
// }

var db *gorm.DB
var cfg config.DBConfig

func InitDB() {
	cfg = config.GetDBConfig()
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBAddr,
		cfg.DBName,
	)
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect mysql")
	} else {
		db = _db
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Envelope{})
		// db.AutoMigrate(&config.CommonConfig{})

		// var checkConfig config.CommonConfig
		// err := db.First(&checkConfig).Error
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	db.Create(config.GetCommonConfig())
		// }
	}
}

func GetDB() *gorm.DB {
	// if db == nil {
	// 	initDB()
	// }
	return db
}
