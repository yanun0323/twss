package delivery

import (
	"main/domain"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	router domain.IRouter
}

func NewHandler(e *echo.Echo, router domain.IRouter) {
	handler := &Handler{
		router: router,
	}

	e.GET("/stock/:id", handler.router.GetStock)
	e.GET("/stock/list", handler.router.GetStockList)
	e.GET("/stock/today", handler.router.GetStocksOfToday)
	e.GET("/stock/top/per", handler.router.GetStocksTopPer)
	e.GET("/stock/top/volume", handler.router.GetStocksTopVolume)
	e.GET("/lastopenday", handler.router.GetLastOpenDay)
}
