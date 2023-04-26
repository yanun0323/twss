package main

import (
	"stocker/internal/app"
	"stocker/pkg/infra"

	"github.com/yanun0323/pkg/logs"
)

func main() {
	l := logs.New("stoker", 2)
	if err := infra.Init("config"); err != nil {
		l.Fatalf("init config , %+v", err)
	}

	app.Run()
}
