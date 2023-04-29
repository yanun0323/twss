package service

import (
	"context"
	"fmt"
	"stocker/internal/domain"
	"stocker/internal/model"
	"stocker/internal/util"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type ConvertRawOption struct {
	Name                 string
	TimeOffset           time.Duration
	StoreMapAfterTask    bool
	GetLastConvertedDate func(Service) (time.Time, error)
	GetRaw               func(Service, time.Time) (domain.Raw, error)
	InsertDataDate       func(Service, time.Time, bool) error
	HandleData           func(Service, context.Context, interface{}) error
}

// ConvertRawTrade 轉換每日盤後交易資料
var ConvertRawTrade = ConvertRawOption{
	Name:              "convert_trade",
	TimeOffset:        -19 * time.Hour, /* turn every 19:00 into 00:00 to convert data after 19:00 every day */
	StoreMapAfterTask: true,
	GetLastConvertedDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetLastTradeDate(svc.Ctx)
	},
	GetRaw: func(svc Service, date time.Time) (domain.Raw, error) {
		return svc.Repo.GetRawTrade(svc.Ctx, date)
	},
	InsertDataDate: func(svc Service, date time.Time, isOpen bool) error {
		return svc.Repo.InsertTradeDate(svc.Ctx, model.TradeDate{
			Date: date,
			Open: isOpen,
		})
	},
	HandleData: func(svc Service, txCtx context.Context, elem interface{}) error {
		trade, ok := elem.(model.Trade)
		if !ok {
			return errors.New("invalid trade type")
		}
		elem, exist := svc.StockMap.Load(trade.ID)
		if !exist {
			stock := trade.CreateStock()
			svc.StockMap.Store(trade.ID, stock)
			elem = stock
		}
		stock := elem.(model.Stock)
		if trade.Date.Before(stock.FirstDate) {
			stock.FirstDate = trade.Date
		}
		if trade.Date.After(stock.LastDate) {
			stock.LastDate = trade.Date
		}

		svc.StockMap.Store(trade.ID, stock)

		if err := svc.Repo.InsertTrade(txCtx, trade); err != nil {
			return errors.Errorf("insert trade, err: %+v", err)
		}

		return nil
	},
}

// ConvertRaw 轉換每日爬蟲資料
func (svc Service) ConvertRaw(opt ConvertRawOption) {
	svc.Log = svc.Log.WithField("service", opt.Name)
	svc.Log.Info("start convert raw...")
	date, err := opt.GetLastConvertedDate(svc)
	if err != nil {
		svc.Log.Errorf("get converted date, %+v", err)
		return
	}
	date = date.Add(24 * time.Hour)
	now := time.Now().Local().Add(opt.TimeOffset)

	wp := util.NewWorkerPool(fmt.Sprintf("convert raw %s", opt.Name), 15)
	wp.Run()

	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		logDate := util.LogDate(date)
		svc.Log.Debugf("start convert raw date: %s", logDate)
		wp.Push(func() {
			err := svc.convertRaw(date, opt)
			if err != nil {
				svc.Log.Errorf("convert raw for %s, err: %+v", logDate, err)
			}
			svc.Log.Infof("finish convert raw date: %s", logDate)
		})
	}

	svc.Log.Info("finish pushing convert task, wait for all task done ...")
	if err := wp.Shutdown(5 * time.Minute); err != nil {
		svc.Log.Errorf("shutdown worker pool, err: %+v", err)
	}

	if opt.StoreMapAfterTask {
		svc.Log.Debug("start store stock map ...")
		if err := svc.storeStockMap(); err != nil {
			svc.Log.Errorf("store stock map, err: %+v", err)
		}
	}

	svc.Log.Info("all trade raw data converted!")
}

func (svc Service) convertRaw(date time.Time, opt ConvertRawOption) error {
	return svc.Repo.Tx(svc.Ctx, func(txCtx context.Context) error {
		raw, err := opt.GetRaw(svc, date)
		if err != nil {
			return errors.Errorf("get raw, err: %+v", err)
		}

		data, err := raw.GetData()
		if err != nil {
			return errors.Errorf("get data, err: %+v", err)
		}

		d, ok := data.(domain.RawData)
		if !ok {
			return errors.New("invalid raw data type")
		}

		isOpen := d.IsOK()
		if !isOpen {
			if err := opt.InsertDataDate(svc, date, false); err != nil {
				return errors.Errorf("insert close data date, err: %+v", err)
			}
			return nil
		}

		elems := d.Parse()
		wg := &sync.WaitGroup{}
		wg.Add(len(elems))
		var repoErr error
		for _, elem := range elems {
			go func(elem interface{}) {
				defer wg.Done()
				if err := opt.HandleData(svc, txCtx, elem); err != nil {
					svc.Log.Errorf("handle data, err: %+v", err)
					repoErr = err
				}
			}(elem)
		}
		wg.Wait()

		if repoErr != nil {
			return repoErr
		}

		if err := opt.InsertDataDate(svc, date, true); err != nil {
			return errors.Errorf("insert open data date, err: %+v", err)
		}

		return nil
	})
}
