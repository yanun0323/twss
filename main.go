package main

import (
	"log"
	"main/config"
	"main/delivery"
	"main/domain"
	"main/model/mode"
	"main/repository"
	"main/service"
	"main/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Init("./config", "config")

	e := echo.New()

	db := repository.ConnectDB()
	repo := repository.NewRepo(db)

	mode.RunTask(config.Mode, mode.Crawl, CrawlData, repo)
	mode.RunTask(config.Mode, mode.Convert, ConvertData, repo)

	if config.Mode != mode.Server {
		return
	}

	service := service.NewService(repo)
	handler := delivery.NewHandler(service)
	delivery.SetRouter(e, handler)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Taiwan Stock Server")
	})

	e.Logger.Fatal(e.Start(":8080"))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	log.Println("shutdown process start")
}

func CrawlData(repo domain.IRepository, checkMode bool) {
	crawler := worker.NewCrawler(repo, checkMode)
	crawler.InitMigrate()
	crawler.Run()
}

func ConvertData(repo domain.IRepository, checkMode bool) {
	converter := worker.NewConverter(repo, checkMode)
	converter.InitMigrate()
	converter.Run()
}
