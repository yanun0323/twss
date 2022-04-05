package domain

import "github.com/labstack/echo/v4"

type IRouter interface {
	GetStock(echo.Context) error
	GetStockList(echo.Context) error
	GetStocksOfToday(echo.Context) error
	GetStocksTopPer(echo.Context) error
	GetStocksTopVolume(echo.Context) error
}
