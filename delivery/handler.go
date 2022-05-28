package delivery

import (
	"errors"
	"main/domain"
	"main/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc service.Service
}

func NewHandler(service service.Service) domain.IHandler {
	return &Handler{svc: service}
}

func (r *Handler) GetStock(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		err := errors.New("Stock id can't be empty.")
		return c.JSON(http.StatusNotFound, ErrorResponse(err))
	}
	stock, err := r.svc.Repo.GetStock(id)
	if err != nil {
		return c.JSON(http.StatusOK, ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, stock)
}
func (r *Handler) GetStockList(c echo.Context) error {
	hash := r.svc.Repo.GetStockHash()
	return c.JSON(http.StatusOK, hash)
}

func (r *Handler) GetStocksOfToday(c echo.Context) error {
	stocks, err := r.svc.Repo.GetStocksToday()
	if err != nil {
		err := errors.New("Failed to get stocks of today.")
		return c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
	return c.JSON(http.StatusOK, stocks)
}
func (r *Handler) GetStocksTopPer(c echo.Context) error {
	// TODO: implement me
	return nil
}
func (r *Handler) GetStocksTopVolume(c echo.Context) error {
	// TODO: implement me
	return nil
}

func (r *Handler) GetLastOpenDay(c echo.Context) error {
	date, err := r.svc.Repo.GetLastOpenDay()
	if err != nil {
		err := errors.New("Failed to get stocks of today.")
		return c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
	return c.JSON(http.StatusOK, date)
}
