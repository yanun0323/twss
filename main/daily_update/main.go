package main

import (
	"main/config"
	"main/domain"
	"main/repository"
	"main/servers"
)

func main() {
	config.Init()

	db := repository.ConnectDB()
	repo := repository.NewRepo(db)

	CrawlData(repo)
	ConvertData(repo)
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
