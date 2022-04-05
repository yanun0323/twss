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

	converter := servers.NewConverter(repo)
	converter.InitMigrate()
	converter.Run()
}

func Convert() {
	main()
}
