package app

import (
	"context"
	"stocker/internal/service"

	"github.com/robfig/cron/v3"
	"github.com/yanun0323/pkg/logs"
)

func CronJob(ctx context.Context, svc service.Service) {
	c := cron.New(cron.WithSeconds())
	l := logs.Get(ctx)
	{
		interval := "TZ=Asia/Taipei 30 0 * * * *"
		_, err := c.AddFunc(interval, func() {
			svc.CrawlRawTrade()
			svc.ConvertRawTrade()
		})
		if err != nil {
			l.Errorf("add trade raw data crawl cron job , %+v", err)
			return
		}
	}

	go c.Run()
	l.Info("processing cron job ...")
}
