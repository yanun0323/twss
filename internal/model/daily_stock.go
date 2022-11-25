package model

import (
	"time"
)

// "證券代號","證券名稱","成交股數","成交金額","開盤價","最高價","最低價","收盤價","漲跌價差","成交筆數"
type DailyStock struct {
	Date                                             time.Time `gorm:"column:date;primaryKey;not null"`
	ID, Name                                         string    `gorm:"-"`
	TradeShare, TradeMoney                           string
	PriceOpen, PriceLowest, PriceHighest, PriceClose string
	TradeGrade, TradeCount                           string
}

func (stock *DailyStock) TableName() string {
	return "stock_" + stock.ID
}
