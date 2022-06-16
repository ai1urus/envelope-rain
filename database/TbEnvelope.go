package database

type Envelope struct {
	Eid         int64 `gorm:"primaryKey"`
	Uid         int64
	Value       int64
	Opened      bool
	Snatch_time int64
}
