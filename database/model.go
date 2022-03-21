package database

type User struct {
	Uid       int64 `gorm:"primaryKey"`
	Amount    int64
	Cur_count int64
}

type Envelope struct {
	Envelope_id int64 `gorm:"primaryKey"`
	Uid         int64
	Value       int64
	Opened      bool
	Snatch_time int64
}
