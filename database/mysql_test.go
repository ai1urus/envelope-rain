package database

import (
	"fmt"
	"testing"
)

func TestDBConnect(t *testing.T) {
	InitDB()
}

func TestDBConfig(t *testing.T) {
	user := &User{
		Uid:       2,
		Amount:    8000,
		Cur_count: 5,
	}
	GetDB().Create(&user)
}

func TestDBSelect(t *testing.T) {
	var result []User
	GetDB().First(&result, 1)
	fmt.Println(result[0].Amount)
}

func TestDBSelectBatch(t *testing.T) {
	var user []User
	result := GetDB().Find(&user)
	fmt.Println(len(user))
	fmt.Println(result.RowsAffected)
}

func TestDBUpdate(t *testing.T) {
	GetDB().Model(&User{}).Where("uid = ?", 1).Update("amount", 100)
	var result []User
	GetDB().First(&result, 1)
	fmt.Println(result[0].Amount)
	// fmt.Println(result[0].Amount)
}

// package main

// import (
// 	"envelope-rain/database"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func main() {
// 	dsn := "root:852196@tcp(175.27.248.24:3306)/envelope?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	db.AutoMigrate(&database.User{})

// 	user := &database.User{
// 		Uid:       2,
// 		Amount:    8000,
// 		Cur_count: 5,
// 	}

// 	db.Create(&user)
// }
