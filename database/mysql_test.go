package database

import (
	"envelope-rain/config"
	"errors"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestDBConnect(t *testing.T) {
	config.InitConfig()
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

func TestDBGetOneNotExist(t *testing.T) {
	config.InitConfig()
	InitDB()
	var envelope Envelope
	result := db.First(&envelope, "1999")
	fmt.Println("FUckyou")
	fmt.Printf("test %v\n", envelope)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("record not found")
	}
}

func TestDBGetBatchNotExist(t *testing.T) {
	config.InitConfig()
	InitDB()
	var envelope []*Envelope
	result := db.Where("uid = ?", 1996).Find(&envelope)
	fmt.Println("FUckyou")
	fmt.Printf("test %v\n", envelope[0])
	// gorm.err
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("record not found")
	}
}

func TestDBUpdateOne(t *testing.T) {
	config.InitConfig()
	InitDB()
	var envelope Envelope
	e := db.First(&Envelope{}, "3").Update("opened", true).Error
	fmt.Println("FUckyou")
	fmt.Printf("test %v\n", envelope)
	if errors.Is(e, gorm.ErrRecordNotFound) {
		fmt.Println("record not found")
	}
}

func TestDBSelectOneRow(t *testing.T) {
	config.InitConfig()
	InitDB()
	var envelope User
	result := db.Select("amount").First(&envelope, 1)
	fmt.Println("FUckyou")
	fmt.Println(envelope)
	fmt.Printf("test %v\n", result)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("record not found")
	}
}

func TestDBUpdateRow(t *testing.T) {
	config.InitConfig()
	InitDB()
	err := UpdateUserValue(1, 10)
	fmt.Println(err)
}

// func TestDBGetUserAmount

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
