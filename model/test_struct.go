package model

import "time"

type TestStruct struct {
	Date  time.Time `gorm:"primaryKey;not null"`
	State bool      `gorm:"index;not null"`
}
