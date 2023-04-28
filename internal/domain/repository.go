package domain

import (
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	DBRepository
	DebugRepository
}

type DBRepository interface {
	ErrRecordNotFound() error

	CheckOpen(time.Time) error
	CheckStock(time.Time) error

	ListTradeRaws(from, to time.Time) ([]model.TradeRaw, error)

	GetLastOpenDate() (time.Time, error)
	GetStockMap() (model.StockMap, error)
	GetStock(id string) (model.Stock, error)
	GetDefaultStartDate() (time.Time, error)
	GetLastTradeRawDate() (time.Time, error)
	GetTradeRaw(time.Time) (model.TradeRaw, error)

	InsertOpen(model.Open) error
	InsertTradeRaw(model.TradeRaw) error
	InsertStockList(model.StockInfo) error
	InsertDailyStockData(model.DailyStock) error
}

type DebugRepository interface {
	Debug() *gorm.DB
}
