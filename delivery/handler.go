package delivery

import (
	"errors"
	"fmt"
	"main/domain"
	"main/model"
	"main/model/compare"
	"main/pkg/response"
	"main/service"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
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
		err := errors.New("stock id can't be empty.")
		return c.JSON(http.StatusNotFound, response.Error(err))
	}
	stock, err := r.svc.Repo.GetStock(id)
	if err != nil {
		return c.JSON(http.StatusOK, response.Error(err))
	}

	return c.JSON(http.StatusOK, stock)
}
func (r *Handler) GetStockList(c echo.Context) error {
	hash := r.svc.Repo.GetStockHash()
	list := make([]model.StockList, 0, len(hash))
	for k, v := range hash {
		list = append(list, model.StockList{StockID: k, StockName: v})
	}
	return c.JSON(http.StatusOK, list)
}

func (r *Handler) GetStocksOfToday(c echo.Context) error {
	stocks, err := r.svc.Repo.GetStocksToday()
	if err != nil {
		err := errors.New("failed to get stocks of today.")
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}
	return c.JSON(http.StatusOK, stocks)
}
func (r *Handler) GetStocksTopPer(c echo.Context) error {
	today, err := r.svc.Repo.GetLastOpenDay()
	if err != nil {
		return err
	}
	stocks, err := r.svc.Repo.GetStocksToday()
	if err != nil {
		return fmt.Errorf("failed to get stocks, %w", err)
	}
	sortableStock := model.NewSortableStock(stocks, today, compare.Per)
	sort.Sort(&sortableStock)
	result := make([]model.Stock, 0, len(sortableStock.Stokes))
	for _, v := range sortableStock.Stokes {
		dec, err := decimal.NewFromString(v.Deals[today].Per)
		if err != nil || dec.IsZero() {
			continue
		}
		result = append(result, v)
	}

	return c.JSON(http.StatusOK, result)
}

func (r *Handler) GetStocksTopVolume(c echo.Context) error {
	today, err := r.svc.Repo.GetLastOpenDay()
	if err != nil {
		return err
	}
	stocks, err := r.svc.Repo.GetStocksToday()
	if err != nil {
		return fmt.Errorf("failed to get stocks, %w", err)
	}
	sortableStock := model.NewSortableStock(stocks, today, compare.Volume)
	sort.Sort(&sortableStock)
	return c.JSON(http.StatusOK, sortableStock.Stokes)
}

func (r *Handler) GetLastOpenDay(c echo.Context) error {
	date, err := r.svc.Repo.GetLastOpenDay()
	if err != nil {
		err := errors.New("failed to get stocks of today.")
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}
	return c.JSON(http.StatusOK, date)
}
