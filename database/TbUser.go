package database

import "gorm.io/gorm"

type User struct {
	Uid       int64 `gorm:"primaryKey"`
	Amount    int64
	Cur_count int64
}

func UpdateUserValue(uid, toAdd int64) (e error) {
	var user User
	e = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("amount").First(&user, uid).Error; err != nil {
			return err
		}
		if err := tx.First(&user, uid).Update("amount", user.Amount+toAdd).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func UpdateUserCount(uid int64) (e error) {
	var user User
	e = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("cur_count").First(&user, uid).Error; err != nil {
			return err
		}
		if err := tx.First(&user, uid).Update("cur_count", user.Cur_count+1).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func GetUserAmount(uid int64) (amount int64, e error) {
	var user User
	e = db.Select("amount").First(&user, uid).Error
	amount = user.Amount
	return
}
