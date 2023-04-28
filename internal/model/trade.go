package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
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
