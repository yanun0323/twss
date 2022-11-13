package service

import (
	"context"
	"stocker/internal/domain"

	"github.com/yanun0323/pkg/logs"
)

type Service struct {
	ctx  context.Context
	l    *logs.Logger
	repo domain.Repository
}

func New(ctx context.Context, repo domain.Repository) Service {
	return Service{
		ctx:  ctx,
		l:    logs.Get(ctx),
		repo: repo,
	}
}
