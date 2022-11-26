package model

import "time"

type Open struct {
	Date   time.Time `gorm:"column:date;primaryKey"`
	IsOpen bool      `gorm:"column:is_open;not null"`
}

func (Open) TableName() string {
	return "open"
}
