package app

import (
	"context"
	"stocker/internal/repository"
	"stocker/internal/service"
)

func Run() {
	ctx := context.Background()
	svc := service.New(ctx, repository.New(ctx))

	svc.CrawlRawDaily()
}
