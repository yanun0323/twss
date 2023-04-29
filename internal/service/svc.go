package service

import (
	"context"
	"stocker/internal/domain"
	"stocker/internal/model"
	"sync"

	"github.com/yanun0323/pkg/logs"
)

type Service struct {
	Ctx      context.Context
	Log      *logs.Logger
	Repo     domain.Repository
	StockMap *sync.Map
}

func New(ctx context.Context, repo domain.Repository) (Service, error) {
	svc := Service{
		Ctx:      ctx,
		Log:      logs.Get(ctx),
		Repo:     repo,
		StockMap: &sync.Map{},
	}

	err := svc.loadStockMap()
	if err != nil {
		return Service{}, err
	}

	return svc, nil
}

func (svc *Service) loadStockMap() error {
	stocks, err := svc.Repo.ListStocks(svc.Ctx)
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		svc.StockMap.Store(stock.ID, stock)
	}
	return nil
}

func (svc *Service) storeStockMap() error {
	return svc.Repo.Tx(svc.Ctx, func(txCtx context.Context) error {
		var err error
		svc.StockMap.Range(func(key, value any) bool {
			stock, ok := value.(model.Stock)
			if !ok {
				return false
			}

			err = svc.Repo.InsertStock(txCtx, stock)
			return err == nil
		})

		return err
	})
}
