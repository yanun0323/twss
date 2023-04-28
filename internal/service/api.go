package service

import (
	"net/http"
	"stocker/internal/util"
	"time"

	"github.com/labstack/echo/v4"
)

func (svc Service) TradeRawAPI(c echo.Context) error {
	d := c.Param("date")
	date, err := time.Parse("20060102", d)
	svc.l.Debug(date)
	if err != nil {
		svc.l.Errorf("[%s] parse date '%s' , %+v", c.RealIP(), d, err)
		return c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid date format", err))
	}

	raw, err := svc.repo.GetRawTrade(date)
	if err != nil {
		svc.l.Errorf("[%s] get daily raw , %+v", c.RealIP(), err)
		return c.JSON(http.StatusInternalServerError, util.NewErrorResponse("internal error", err))
	}

	svc.l.Infof("[%s] get daily raw succeed", c.RealIP())
	return c.JSON(http.StatusOK, util.NewDataResponse("get daily raw data succeed", string(raw.Body)))
}
