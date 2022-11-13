package domain

import (
	"stocker/internal/model"
	"time"
)

type Repository interface {
	DBRepository
}

type DBRepository interface {
	ListDailyRaws(from, to time.Time) ([]model.DailyRaw, error)
	ListAllDailyRaws() ([]model.DailyRaw, error)
	GetLastDailyRaw() (model.DailyRaw, error)
	InsertDailyRaw(raw model.DailyRaw) error

	InsertDailyStock(stock model.DailyStock) error
}
