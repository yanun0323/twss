package model

import "time"

type TradeDate struct {
	Date   time.Time `gorm:"column:date;primaryKey"`
	IsOpen bool      `gorm:"column:is_open;not null"`
}

func (TradeDate) TableName() string {
	return "trade_date"
}
