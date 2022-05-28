package delivery

import (
	"main/domain"

	"github.com/labstack/echo/v4"
)

func SetRouter(e *echo.Echo, handler domain.IHandler) {
	api := e.Group("/api")
	api.GET("/stock/:id", handler.GetStock)
	api.GET("/stock/list", handler.GetStockList)
	api.GET("/stock/today", handler.GetStocksOfToday)
	api.GET("/stock/top/per", handler.GetStocksTopPer)
	api.GET("/stock/top/volume", handler.GetStocksTopVolume)
	api.GET("/lastopenday", handler.GetLastOpenDay)
}
