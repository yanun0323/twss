package service

import (
	"net/http"
	"stocker/internal/util"
	"time"

	"github.com/labstack/echo/v4"
)

func (svc Service) RawDailyAPI(c echo.Context) error {
	d := c.QueryParam("date")
	date, err := time.Parse("20060102", d)
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
	return c.JSON(http.StatusOK, util.NewDataResponse("Get daily raw data succeed", struct {
		Date string `json:"date"`
		Body string `json:"body"`
	}{
		Date: raw.Date.Format("2006-01-02"),
		Body: string(raw.Body),
	}))
}
