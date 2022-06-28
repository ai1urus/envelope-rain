package database

import "gorm.io/gorm"

type Envelope struct {
	Eid         int64 `gorm:"primaryKey"`
	Uid         int64
	Value       int64
	Opened      bool
	Snatch_time int64
}

func CreateEnvelope(envelope Envelope) (e error) {
	e = db.Create(&envelope).Error
	return
}

func GetEnvelopeByEid(eid int64) (envelope Envelope, e error) {
	e = db.First(&envelope, eid).Error
	return
}

func GetEnvelopeByUid(uid int64) (envelopes []*Envelope, e error) {
	db.Where("uid = ?", uid).Find(&envelopes)
	if len(envelopes) == 0 {
		e = gorm.ErrRecordNotFound
	}
	return
}

func SetEnvelopeOpen(eid int64) (e error) {
	e = db.First(&Envelope{}, eid).Update("opened", true).Error
	return
}
