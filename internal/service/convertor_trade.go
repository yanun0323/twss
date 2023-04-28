package service

import (
	"context"
	"stocker/internal/model"
	"time"
)

func (svc Service) ConvertRawTrade() {
	svc.l.Info("start trade raw data convert ...")
	date, err := svc.repo.GetLastTradeDate(svc.ctx)
	if err != nil {
		svc.l.Errorf("get last open date , %+v", err)
		return
	}
	date = date.Add(24 * time.Hour)
	now := time.Now()

	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		err := svc.convertRawTrade(date)
		if err != nil {
			svc.l.Errorf("convert raw trade, err: %+v", err)
			return
		}
	}

	svc.l.Info("all trade raw data converted!")
}

func (svc Service) convertRawTrade(date time.Time) error {
	return svc.repo.Tx(svc.ctx, func(txCtx context.Context) error {
		raw, err := svc.repo.GetRawTrade(txCtx, date)
		if err != nil {
			return err
		}

		data, err := raw.GetData()
		if err != nil {
			return err
		}

		if data.Stat != "OK" || len(data.TradeData()) == 0 {
			return svc.repo.InsertTradeDate(txCtx, model.TradeDate{
				Date:   date,
				IsOpen: false,
			})
		}

		trades := data.ParseTrade()
		for _, trade := range trades {
			elem, exist := svc.stockMap.Load(trade.ID)
			if !exist {
				svc.stockMap.Store(trade.ID, trade.CreateStock())
			}
			stock := elem.(model.Stock)
			stock.LastDate = date
			svc.stockMap.Store(trade.ID, stock)

			if err := svc.repo.InsertTrade(txCtx, trade); err != nil {
				return err
			}
		}

		return svc.repo.InsertTradeDate(txCtx, model.TradeDate{
			Date:   date,
			IsOpen: true,
		})
	})
}
