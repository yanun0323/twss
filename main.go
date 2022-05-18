package main

import (
	"log"
	"main/config"
	"main/delivery"
	"main/domain"
	"main/repository"
	"main/servers"
	"main/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Init()

	db := repository.ConnectDB()
	repo := repository.NewRepo(db)

	CrawlData(repo)
	ConvertData(repo)

	router := usecase.NewRouter(repo)

	e := echo.New()
	delivery.NewHandler(e, router)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Taiwan Stock Server")
	})

	e.Logger.Fatal(e.Start(":8080"))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	log.Println("shutdown process start")
}

func CrawlData(repo domain.IRepository) {
	crawler := servers.NewCrawler(repo)
	crawler.InitMigrate()
	crawler.Run()
}

func ConvertData(repo domain.IRepository) {
	converter := servers.NewConverter(repo)
	converter.InitMigrate()
	converter.Run()
}
