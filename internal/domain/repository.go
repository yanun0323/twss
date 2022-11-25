package domain

import (
	"stocker/internal/model"
	"time"
)

type Repository interface {
	DBRepository
}

type DBRepository interface {
	ListAllDailyRaws() ([]model.DailyRaw, error)
	ListDailyRaws(from, to time.Time) ([]model.DailyRaw, error)

	GetLastOpenDate() (time.Time, error)
	GetStockMap() (model.StockMap, error)
	GetLastDailyRawDate() (time.Time, error)
	GetDailyRaw(time.Time) (model.DailyRaw, error)

	InsertOpen(model.Open) error
	InsertDailyRaw(model.DailyRaw) error
	InsertStockList(model.StockInfo) error
	InsertDailyStock(model.DailyStock) error
}
