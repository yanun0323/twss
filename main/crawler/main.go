package main

import (
	"main/config"
	"main/dao"
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

func Crawl() {
	main()
}
