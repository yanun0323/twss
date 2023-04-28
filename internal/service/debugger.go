package service

import (
	"time"
)

func (svc Service) Debug() {
	svc.l.Info("start debugging ...")

	if false { /* to avoid unused warning */
		svc.debugTradeRawData("20221125")
		// svc.debugRefactorStockList()
	}
}

func (svc Service) debugTradeRawData(dateStr string) {
	date, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	raw, err := svc.repo.GetRawTrade(svc.ctx, date)
	if err != nil {
		svc.l.Errorf("get trade raw , %+v", err)
		return
	}
	d, err := raw.GetData()
	if err != nil {
		svc.l.Errorf("get data , %+v", err)
		return
	}

	svc.l.Debugf("%+v", d.Date)
	svc.l.Debugf("%+v", d.Fields8)
	svc.l.Debugf("%+v", d.TradeData()[0])
	svc.l.Debugf("%+v", d.ParseTrade()[0])
}

// func (svc Service) debugRefactorStockList() {
// 	db := svc.repo.Debug(svc.ctx)

// 	stockMap, err := svc.repo.GetStockMap(svc.ctx)
// 	if err != nil {
// 		svc.l.Errorf("get stock map , %+v", err)
// 		return
// 	}

// 	lastOpen := model.TradeDate{}
// 	if err := db.Where("is_open = ?", true).Last(&lastOpen).Error; err != nil {
// 		svc.l.Errorf("get last open date , %+v", err)
// 		return
// 	}
// 	for id, name := range stockMap {
// 		table := model.Trade{ID: id}.GetTableName()
// 		start, last := model.Trade{}, model.Trade{}
// 		if err := db.Table(table).First(&start).Error; err != nil {
// 			svc.l.Errorf("%s, get stock first date , %+v", id, err)
// 			return
// 		}
// 		if err := db.Table(table).Last(&last).Error; err != nil {
// 			svc.l.Errorf("%s, get stock last date , %+v", id, err)
// 			return
// 		}
// 		unit := model.Stock{
// 			ID:        id,
// 			Name:      name,
// 			FirstDate: start.Date,
// 			LastDate:  last.Date,
// 			Unable:    !last.Date.Equal(lastOpen.Date),
// 		}
// 		if last.Date.Equal(lastOpen.Date) {
// 			svc.l.Debugf("%s - %s Enable", id, name)
// 		}
// 		if err := db.Model(&unit).Updates(unit).Error; err != nil {
// 			svc.l.Errorf("update stock list , %+v", err)
// 			return
// 		}
// 	}
// }
