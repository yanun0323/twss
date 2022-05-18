package main

import (
	"main/config"
	dao "main/repository"
	"main/servers"
)

func main() {
	config.Init()

	db := dao.ConnectDB()
	repo := dao.NewRepo(db)

	crawler := servers.NewCrawler(repo)
	crawler.InitMigrate()
	crawler.Run()
}
