package service

import (
	"context"
	"stocker/internal/domain"
	"stocker/internal/model"
	"sync"

	"github.com/yanun0323/pkg/logs"
)

type Service struct {
	ctx      context.Context
	l        *logs.Logger
	repo     domain.Repository
	stockMap *sync.Map
}

func New(ctx context.Context, repo domain.Repository) (Service, error) {
	svc := Service{
		ctx:      ctx,
		l:        logs.Get(ctx),
		repo:     repo,
		stockMap: &sync.Map{},
	}

	err := svc.loadStockMap()
	if err != nil {
		return Service{}, err
	}

	return svc, nil
}

func (svc *Service) loadStockMap() error {
	stocks, err := svc.repo.ListStocks(svc.ctx)
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		svc.stockMap.Store(stock.ID, stock)
	}
	return nil
}

func (svc *Service) storeStockMap() error {
	return svc.repo.Tx(svc.ctx, func(txCtx context.Context) error {
		var err error
		svc.stockMap.Range(func(key, value any) bool {
			stock, ok := value.(model.Stock)
			if !ok {
				return false
			}

			err = svc.repo.InsertStock(txCtx, stock)
			return err == nil
		})

		return err
	})
}
