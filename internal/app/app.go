package app

import (
	"context"
	"os"
	"os/signal"
	"stocker/internal/repository"
	"stocker/internal/service"
	"syscall"

	"github.com/yanun0323/pkg/logs"
)

func Run() {
	ctx := context.Background()
	svc := service.New(ctx, repository.New(ctx))
	CronJob(ctx, svc)
	APIServer(ctx, svc)

	svc.CrawlDailyRawData()
	// svc.ConvertDailyRawData()

	/* Graceful shutdown */
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	logs.Get(ctx).Info("shutdown process start")
}
