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
	TradeRepository
	FinanceRepository

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

	ListRawFinance(ctx context.Context, from, to time.Time) ([]model.RawFinance, error)
	GetRawFinanceDate(ctx context.Context, begin bool) (time.Time, error)
	GetRawFinance(context.Context, time.Time) (model.RawFinance, error)
	InsertRawFinance(context.Context, model.RawFinance) error
}

type TradeRepository interface {
	IsTradeExist(context.Context, time.Time) (bool, error)
	ListTrade(ctx context.Context, id string, from, to time.Time) ([]model.Trade, error)
	GetTrade(ctx context.Context, id string, date time.Time) (model.Trade, error)
	InsertTrade(context.Context, model.Trade) error

	IsTradeDateExist(context.Context, time.Time) (bool, error)
	GetTradeDate(context.Context, time.Time) (model.TradeDate, error)
	GetLastTradeDate(context.Context) (time.Time, error)
	InsertTradeDate(context.Context, model.TradeDate) error
}

type FinanceRepository interface {
	IsFinanceExist(context.Context, time.Time) (bool, error)
	InsertFinance(context.Context, model.Finance) error

	IsFinanceDateExist(context.Context, time.Time) (bool, error)
	GetFinanceDate(context.Context, time.Time) (model.FinanceDate, error)
	GetLastFinanceDate(context.Context) (time.Time, error)
	InsertFinanceDate(context.Context, model.FinanceDate) error
}

type DebugRepository interface {
	Debug(context.Context) *gorm.DB
}
