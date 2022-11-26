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

	ListDailyRaws(from, to time.Time) ([]model.DailyRaw, error)

	GetLastOpenDate() (time.Time, error)
	GetStockMap() (model.StockMap, error)
	GetStock(id string) (model.Stock, error)
	GetDefaultStartDate() (time.Time, error)
	GetLastDailyRawDate() (time.Time, error)
	GetDailyRaw(time.Time) (model.DailyRaw, error)

	InsertOpen(model.Open) error
	InsertDailyRaw(model.DailyRaw) error
	InsertStockList(model.StockInfo) error
	InsertDailyStock(model.DailyStock) error
}

type DebugRepository interface {
	Debug() *gorm.DB
}
