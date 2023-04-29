package model

import (
	"time"
)

type Stock struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	FirstDate time.Time `gorm:"column:first_date;not null"`
	LastDate  time.Time `gorm:"column:last_date;not null"`
	Unable    bool      `gorm:"column:unable;not null"`
}

func (Stock) TableName() string {
	return "stocks"
}
