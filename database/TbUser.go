package database

type User struct {
	Uid       int64 `gorm:"primaryKey"`
	Amount    int64
	Cur_count int64
}
