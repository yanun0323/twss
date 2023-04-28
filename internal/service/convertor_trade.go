package service

import (
	"context"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"

	"github.com/pkg/errors"
)

func (svc Service) ConvertRawTrade() {
	svc.l.Info("start trade raw data convert ...")
	date, err := svc.repo.GetLastTradeDate(svc.ctx)
	if err != nil {
		svc.l.Errorf("get last open date , %+v", err)
		return
	}
	date = date.Add(24 * time.Hour)
	now := time.Now().Local().Add(-19 * time.Hour) /* turn every 19:00 into 00:00 to convert data after 19:00 every day */

	wp := util.NewWorkerPool("convert raw trade", 15)
	wp.Run()

	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		logDate := util.LogDate(date)
		svc.l.Debugf("start convert trade date: %s", logDate)
		wp.Push(func() {
			err := svc.convertRawTrade(date)
			if err != nil {
				svc.l.Errorf("convert raw trade for %s, err: %+v", logDate, err)
			}
			svc.l.Infof("finish convert trade date: %s", logDate)
		})
	}

	svc.l.Info("finish pushing convert task, wait for all task done ...")
	if err := wp.Shutdown(5 * time.Minute); err != nil {
		svc.l.Errorf("shutdown worker pool, err: %+v", err)
	}

	svc.l.Info("start store stock map ...")
	if err := svc.storeStockMap(); err != nil {
		svc.l.Errorf("store stock map, err: %+v", err)
	}

	svc.l.Info("all trade raw data converted!")
}

func (svc Service) convertRawTrade(date time.Time) error {
	return svc.repo.Tx(svc.ctx, func(txCtx context.Context) error {
		raw, err := svc.repo.GetRawTrade(txCtx, date)
		if err != nil {
			return errors.Errorf("get raw trade, err: %+v", err)
		}

		data, err := raw.GetData()
		if err != nil {
			return errors.Errorf("get data, err: %+v", err)
		}

		if data.Stat != "OK" || len(data.TradeData()) == 0 {
			if err := svc.repo.InsertTradeDate(txCtx, model.TradeDate{
				Date:   date,
				IsOpen: false,
			}); err != nil {
				return errors.Errorf("insert close trade date, err: %+v", err)
			}
			return nil
		}

		trades := data.ParseTrade()
		for _, trade := range trades {
			elem, exist := svc.stockMap.Load(trade.ID)
			if !exist {
				stock := trade.CreateStock()
				svc.stockMap.Store(trade.ID, stock)
				elem = stock
			}
			stock := elem.(model.Stock)
			if date.Before(stock.FirstDate) {
				stock.FirstDate = date
			}
			if date.After(stock.LastDate) {
				stock.LastDate = date
			}

			svc.stockMap.Store(trade.ID, stock)

			if err := svc.repo.InsertTrade(txCtx, trade); err != nil {
				return errors.Errorf("insert trade, err: %+v", err)
			}
		}

		if err := svc.repo.InsertTradeDate(txCtx, model.TradeDate{
			Date:   date,
			IsOpen: true,
		}); err != nil {
			return errors.Errorf("insert open trade date, err: %+v", err)
		}

		return nil
	})
}
