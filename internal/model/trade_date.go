package model

import "time"

type TradeDate struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Open bool      `gorm:"column:is_open;not null"`
}

func (td TradeDate) IsOpen() bool {
	return td.Open
}

func (TradeDate) TableName() string {
	return "trade_date"
}
