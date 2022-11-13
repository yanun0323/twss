package main

import (
	"stocker/internal/app"

	"github.com/yanun0323/pkg/config"
	"github.com/yanun0323/pkg/logs"
)

func main() {
	l := logs.New("stoker", 2)
	if err := config.Init("config"); err != nil {
		l.Fatalf("init config failed, %+v", err)
	}

	app.Run()
}
