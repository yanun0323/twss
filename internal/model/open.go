package model

import "time"

type Open struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Open bool      `gorm:"column:open;not null"`
}

func (Open) TableName() string {
	return "open"
}
