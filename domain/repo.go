package domain

import (
	"main/model"
	"time"
)

type IRepository interface {
	GetCrawlableDate() time.Time
	Insert(interface{}) error
	Create(string, interface{}) error
	AutoMigrate(...interface{}) error
	Migrate(string, interface{}) error
	GetConvertableDate() (time.Time, bool)
	GetRaw(time.Time) ([]byte, error)
	GetStockHash() map[string]string
	GetLastOpenDay() (time.Time, error)

	GetStock(string) (model.Stock, error)
	GetStocksToday() ([]model.Stock, error)
}
