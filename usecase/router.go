package usecase

import (
	"main/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Router struct {
	repo domain.IRepository
}

func NewRouter(repo domain.IRepository) domain.IRouter {
	return &Router{repo: repo}
}

func (r *Router) GetStock(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "Stock id can't be empty.")
	}
	stock, err := r.repo.GetStock(id)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, stock)
}
func (r *Router) GetStockList(c echo.Context) error {
	hash := r.repo.GetStockHash()
	return c.JSON(http.StatusOK, hash)
}

func (r *Router) GetStocksOfToday(c echo.Context) error {
	stocks, err := r.repo.GetStocksToday()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get stocks of today.")
	}
	return c.JSON(http.StatusOK, stocks)
}
func (r *Router) GetStocksTopPer(c echo.Context) error {
	// TODO: implement me
	return nil
}
func (r *Router) GetStocksTopVolume(c echo.Context) error {
	// TODO: implement me
	return nil
}

func (r *Router) GetLastOpenDay(c echo.Context) error {
	date, err := r.repo.GetLastOpenDay()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get last open day.")
	}
	return c.JSON(http.StatusOK, date)
}
