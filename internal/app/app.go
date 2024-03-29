package app

import (
	"context"
	"log"
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
	repo, err := repository.New(ctx)
	if err != nil {
		log.Fatalf("init repository, err: %+v", err)
	}

	svc, err := service.New(ctx, repo)
	if err != nil {
		log.Fatalf("init service, err: %+v", err)
	}

	mode := strings.ToLower(os.Getenv("MODE"))
	logs.Get(ctx).Infof("MODE: %s", mode)
	switch mode {
	case "check":
		RunCheck(svc)
	case "repair":
		RunRepair(svc)
	case "job":
		RunJobOnce(svc)
	case "debug":
		RunDebug(svc)
	default:
		RunJobOnce(svc)
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
	svc.CheckRaw(false, service.CheckRawTrade)
	svc.CheckRaw(false, service.CheckRawFinance)
	svc.CheckData(false, service.CheckTrade)
	svc.CheckData(false, service.CheckFinance)
}

func RunRepair(svc service.Service) {
	svc.CheckRaw(true, service.CheckRawTrade)
	svc.CheckRaw(true, service.CheckRawFinance)
	svc.CheckData(true, service.CheckTrade)
	svc.CheckData(true, service.CheckFinance)
}

func RunJobOnce(svc service.Service) {
	svc.CrawlRaw(service.CrawlFinance)
	svc.CrawlRaw(service.CrawlTrade)
	svc.ConvertRaw(service.ConvertRawTrade)
	svc.ConvertRaw(service.ConvertRawFinance)
}

func RunDebug(svc service.Service) {
	svc.Debug()
}
