package main

import (
	"log"
	"main/config"
	"main/dao"
	"main/delivery"
	"main/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Init()
	db := dao.ConnectDB()
	repo := dao.NewRepo(db)
	router := usecase.NewRouter(repo)

	e := echo.New()
	delivery.NewHandler(e, router)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "TWSS")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	log.Println("shutdown process start")
}
