package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stock struct {
	Date         time.Time       `gorm:"column:date;primaryKey" json:"date"`
	ID           string          `gorm:"-" json:"id,omitempty"`
	Name         string          `gorm:"-" json:"name,omitempty"`
	TradeShare   decimal.Decimal `gorm:"not null" json:"tradeShare"`
	TradeCount   decimal.Decimal `gorm:"not null" json:"tradeCount"`
	TradeMoney   decimal.Decimal `gorm:"not null" json:"tradeMoney"`
	PriceOpen    decimal.Decimal `gorm:"not null" json:"priceOpen"`
	PriceHighest decimal.Decimal `gorm:"not null" json:"priceHighest"`
	PriceLowest  decimal.Decimal `gorm:"not null" json:"priceLowest"`
	PriceClose   decimal.Decimal `gorm:"not null" json:"priceClose"`
	TradeSymbol  string          `gorm:"not null" json:"tradeSymbol"`
	TradeGrade   decimal.Decimal `gorm:"not null" json:"tradeGrade"`
	Percentage   decimal.Decimal `gorm:"not null" json:"percentage"`
	Limit        bool            `gorm:"not null" json:"limit"`
}

func (stock Stock) GetTableName() string {
	return "stock_" + stock.ID
}

// type Stock struct {
// 	ID        string       `json:"id"`
// 	Name      string       `json:"name"`
// 	FirstDate time.Time    `json:"firstDate"`
// 	LastDate  time.Time    `json:"lastDate"`
// 	Unable    bool         `json:"unable"`
// 	Trading   []DailyStock `json:"trading"`
// }
