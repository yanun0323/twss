package model

import "time"

type OpenDays struct {
	Date  time.Time `gorm:"primaryKey;not null"`
	State bool      `gorm:"index;not null"`
}

func (c OpenDays) TableName() string {
	return "open_days"
}
