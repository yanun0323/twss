package domain

import (
	"context"
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CommonRepository
	StockRepository
	RawRepository
	TradeDateRepository
	TradeRepository

	DebugRepository
}

type CommonRepository interface {
	Tx(ctx context.Context, fc func(txCtx context.Context) error) error
	ErrNotFound() error
	GetDefaultStartDate() (time.Time, error)
}

type StockRepository interface {
	ListStocks(context.Context) ([]model.Stock, error)
	InsertStock(context.Context, model.Stock) error
}

type RawRepository interface {
	ListRawTrade(ctx context.Context, from, to time.Time) ([]model.RawTrade, error)
	GetRawTradeDate(ctx context.Context, begin bool) (time.Time, error)
	GetRawTrade(context.Context, time.Time) (model.RawTrade, error)
	InsertRawTrade(context.Context, model.RawTrade) error

	ListRawEps(ctx context.Context, from, to time.Time) ([]model.RawEps, error)
	GetRawEpsDate(ctx context.Context, begin bool) (time.Time, error)
	GetRawEps(context.Context, time.Time) (model.RawEps, error)
	InsertRawEps(context.Context, model.RawEps) error
}

type TradeRepository interface {
	IsTradeExist(context.Context, time.Time) (bool, error)
	InsertTrade(context.Context, model.Trade) error
}

type TradeDateRepository interface {
	IsTradeDateExist(context.Context, time.Time) (bool, error)
	GetTradeDate(context.Context, time.Time) (model.TradeDate, error)
	GetLastTradeDate(context.Context) (time.Time, error)
	InsertTradeDate(context.Context, model.TradeDate) error
}

type DebugRepository interface {
	Debug(context.Context) *gorm.DB
}
