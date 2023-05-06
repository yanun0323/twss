package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trade 交易資料
type Trade struct {
	Date         time.Time       `gorm:"column:date;primaryKey" json:"date"`
	ID           string          `gorm:"-" json:"id,omitempty"`
	Name         string          `gorm:"-" json:"name,omitempty"`
	TradeShare   decimal.Decimal `gorm:"not null" json:"tradeShare"`   /* 成交股數 */
	TradeCount   decimal.Decimal `gorm:"not null" json:"tradeCount"`   /* 成交筆數 */
	TradeMoney   decimal.Decimal `gorm:"not null" json:"tradeMoney"`   /* 成交金額 */
	PriceOpen    decimal.Decimal `gorm:"not null" json:"priceOpen"`    /* 開盤價 */
	PriceHighest decimal.Decimal `gorm:"not null" json:"priceHighest"` /* 最高價 */
	PriceLowest  decimal.Decimal `gorm:"not null" json:"priceLowest"`  /* 最低價 */
	PriceClose   decimal.Decimal `gorm:"not null" json:"priceClose"`   /* 收盤價 */
	TradeSymbol  string          `gorm:"not null" json:"tradeSymbol"`  /* 漲跌前綴符號 */
	TradeGrade   decimal.Decimal `gorm:"not null" json:"tradeGrade"`   /* 漲跌價差 */
	Percentage   decimal.Decimal `gorm:"not null" json:"percentage"`   /* 漲跌百分比 */
	Limit        bool            `gorm:"not null" json:"limit"`        /* 是否漲跌停 */
}

func (trade Trade) GetTableName() string {
	return "trade_" + trade.ID
}

func (trade Trade) CreateStock() Stock {
	return Stock{
		ID:        trade.ID,
		Name:      trade.Name,
		FirstDate: trade.Date,
		Unable:    false,
	}
}

// TradeDate 交易日期
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
