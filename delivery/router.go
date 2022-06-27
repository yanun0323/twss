package delivery

import (
	"main/domain"
	"main/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func SetRouter(e *echo.Echo, handler domain.IHandler) {
	api := e.Group("/api", middleware.TokenAuth)
	api.GET("/stock/:id", handler.GetStock)
	api.GET("/stock/list", handler.GetStockList)
	api.GET("/stock/today", handler.GetStocksOfToday)
	api.GET("/lastopenday", handler.GetLastOpenDay)
}
