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

func initDB() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetConfig().GetString("mysql.username"),
		config.GetConfig().GetString("mysql.password"),
		config.GetConfig().GetString("mysql.address"),
		config.GetConfig().GetString("mysql.dbname"),
	)
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect mysql")
	} else {
		db = _db
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Envelope{})
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		initDB()
	}
	return db
}
