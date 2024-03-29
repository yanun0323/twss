package app

import (
	"context"
	"net/http"
	"stocker/internal/service"
	"stocker/internal/util"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

func APIServer(ctx context.Context, svc service.Service) {
	l := logs.Get(ctx)
	e := echo.New()

	api := e.Group("/api", DefaultMiddleware(ctx)...)
	api.GET("/healthz", func(c echo.Context) error {
		l.Infof("[%s] health check succeed", c.RealIP())
		return c.JSON(http.StatusOK, util.NewMsgResponse("health check succeed :)"))
	})
	api.GET("/raw/trade/:date", svc.TradeRawAPI)

	port := ":" + viper.GetString("server.port")
	go e.Start(port)
	<-time.After(time.Second)
}
