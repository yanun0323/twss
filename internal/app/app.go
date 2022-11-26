package app

import (
	"context"
	"os"
	"os/signal"
	"stocker/internal/repository"
	"stocker/internal/service"
	"strings"
	"syscall"

	"github.com/yanun0323/pkg/logs"
)

func Run() {
	ctx := context.Background()
	svc := service.New(ctx, repository.New(ctx))

	mode := strings.ToLower(os.Getenv("MODE"))
	logs.Get(ctx).Infof("MODE: %s", mode)
	switch mode {
	case "check":
		RunCheck(svc)
	case "job":
		RunJob(svc)
	case "debug":
		RunDebug(svc)
	default:
		RunJob(svc)
		RunServer(ctx, svc)
	}
}

func RunServer(ctx context.Context, svc service.Service) {
	CronJob(ctx, svc)
	APIServer(ctx, svc)
	/* Graceful shutdown */
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	logs.Get(ctx).Warn("processing shutdown ...")
}

func RunCheck(svc service.Service) {
	svc.CheckDailyRaw()
	svc.CheckConverter()
}

func RunJob(svc service.Service) {
	svc.CrawlDailyRawData()
	svc.ConvertDailyRawData()
}

func RunDebug(svc service.Service) {
	svc.Debug()
}
