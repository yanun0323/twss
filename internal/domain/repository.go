package domain

import (
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CommonRepository
	OpenRepository
	RawRepository
	StockRepository

	DebugRepository
}

type CommonRepository interface {
	ErrRecordNotFound() error

	GetDefaultStartDate() (time.Time, error)

	InsertStockList(model.StockListUnit) error
}

type OpenRepository interface {
	CheckOpen(time.Time) error
	GetLastOpenDate() (time.Time, error)
	InsertOpen(model.Open) error
}

type RawRepository interface {
	ListRawTrades(from, to time.Time) ([]model.RawTrade, error)
	GetLastRawTradeDate() (time.Time, error)
	GetRawTrade(time.Time) (model.RawTrade, error)
	InsertRawTrade(model.RawTrade) error

	// ListEpsRaws(from, to time.Time) ([]model.EpsRaw, error)
}

type StockRepository interface {
	CheckStock(time.Time) error
	GetStockMap() (model.StockMap, error)
	InsertStock(model.Stock) error
}

type DebugRepository interface {
	Debug() *gorm.DB
}
