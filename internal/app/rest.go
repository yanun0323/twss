package app

import (
	"context"
	"net/http"
	"stocker/internal/service"
	"stocker/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

func APIServer(ctx context.Context, svc service.Service) {
	l := logs.Get(ctx)
	e := echo.New()

	public := e.Group("/public")
	public.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, util.NewMsgResponse("OK"))
	})

	port := ":" + viper.GetString("server.port")
	go e.Start(port)
	l.Infof("start api server at port %s", port)
}
