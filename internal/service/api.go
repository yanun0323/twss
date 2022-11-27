package service

import (
	"errors"
	"net/http"
	"stocker/internal/util"
	"time"

	"github.com/labstack/echo/v4"
)

func (svc Service) RawDailyAPI(c echo.Context) error {
	d := c.Param("date")
	date, err := time.Parse("20060102", d)
	svc.l.Debug(date)
	if err != nil {
		svc.l.Errorf("[%s] parse date '%s' failed, %+v", c.RealIP(), d, err)
		return c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid date format", err))
	}

	raw, err := svc.repo.GetDailyRaw(date)
	if err != nil {
		svc.l.Errorf("[%s] get daily raw failed, %+v", c.RealIP(), err)
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("internal error", err))
	}

	svc.l.Infof("[%s] get daily raw succeed", c.RealIP())
	return c.JSON(http.StatusOK, util.NewDataResponse("Get daily raw data succeed", string(raw.Body)))
}

func (svc Service) StockDailyAPI(c echo.Context) error {
	d := c.Param("date")
	date, err := time.Parse("20060102", d)
	if err != nil {
		svc.l.Errorf("[%s] parse date '%s' failed, %+v", c.RealIP(), d, err)
		return c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid date format", err))
	}

	raw, err := svc.repo.GetDailyRaw(date)
	if err != nil {
		svc.l.Errorf("[%s] get daily stock failed, %+v", c.RealIP(), err)
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("internal error", err))
	}

	data, err := raw.GetData()
	if err != nil {
		svc.l.Errorf("[%s] get raw data failed, %+v", c.RealIP(), err)
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("internal error", err))
	}

	data.ParseStockList()

	svc.l.Infof("[%s] get daily stock succeed", c.RealIP())
	return c.JSON(http.StatusOK, util.NewDataResponse("Get daily stock succeed", data))
}

func (svc Service) StockAPI(c echo.Context) error {
	svc.l.Debug(c.ParamNames())
	id := c.Param("id")
	if len(id) == 0 {
		svc.l.Errorf("[%s] empty stock id", c.RealIP())
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("empty stock id"))
	}
	stock, err := svc.repo.GetStock(id)
	if errors.Is(svc.repo.ErrRecordNotFound(), err) {
		svc.l.Errorf("[%s] invalid stock id, %+v", c.RealIP(), err)
		return c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid stock id"))
	}
	if err != nil {
		svc.l.Errorf("[%s] get stock failed, %+v", c.RealIP(), err)
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("internal error", err))
	}

	svc.l.Infof("[%s] get stock succeed", c.RealIP())
	return c.JSON(http.StatusOK, util.NewDataResponse("Get stock succeed", stock))
}
