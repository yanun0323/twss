package main

import (
	"main/config"
	"main/repository"
	"main/servers"
)

func main() {
	config.Init()

	db := repository.ConnectDB()
	repo := repository.NewRepo(db)

	converter := servers.NewConverter(repo)
	converter.InitMigrate()
	converter.Run()
}
