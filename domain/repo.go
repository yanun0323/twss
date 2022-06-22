package domain

import (
	"main/model"
	"time"
)

type IRepository interface {
	GetCrawlableDate(checkMode bool) time.Time

	GetConvertibleDate(checkMode bool) time.Time
	GetRaw(date time.Time) ([]byte, error)
	GetLastOpenDay() (time.Time, error)

	GetStock(id string) (model.Stock, error)
	GetStocksToday() ([]model.Stock, error)
	GetStockHash() map[string]string

	AutoMigrate(p ...interface{}) error
	Migrate(table string, p interface{}) error
	Insert(obj interface{}) error
	InsertWithTableName(table string, obj interface{}) error
}
