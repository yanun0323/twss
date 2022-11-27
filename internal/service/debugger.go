package service

import (
	"stocker/internal/model"
	"time"
)

func (svc Service) Debug() {
	svc.l.Info("start debugging ...")

	if false { /* to avoid unused warning */
		svc.debugDailyRawData("20221125")
		svc.debugRefactorStockList()
	}
}

func (svc Service) debugDailyRawData(dateStr string) {
	date, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	raw, err := svc.repo.GetDailyRaw(date)
	if err != nil {
		svc.l.Errorf("get daily raw failed, %+v", err)
		return
	}
	d, err := raw.GetData()
	if err != nil {
		svc.l.Errorf("get data failed, %+v", err)
		return
	}

	svc.l.Debugf("%+v", d.RawDate)
	svc.l.Debugf("%+v", d.Fields8)
	svc.l.Debugf("%+v", d.Data()[0])
	svc.l.Debugf("%+v", d.ParseStock()[0])
}

func (svc Service) debugRefactorStockList() {
	db := svc.repo.Debug()

	stockMap, err := svc.repo.GetStockMap()
	if err != nil {
		svc.l.Errorf("get stock map failed, %+v", err)
		return
	}

	lastOpen := model.Open{}
	if err := db.Where("is_open = ?", true).Last(&lastOpen).Error; err != nil {
		svc.l.Errorf("get last open date failed, %+v", err)
		return
	}
	for id, name := range stockMap {
		table := model.DailyStock{ID: id}.GetTableName()
		start, last := model.DailyStock{}, model.DailyStock{}
		if err := db.Table(table).First(&start).Error; err != nil {
			svc.l.Errorf("%s, get stock first date failed, %+v", id, err)
			return
		}
		if err := db.Table(table).Last(&last).Error; err != nil {
			svc.l.Errorf("%s, get stock last date failed, %+v", id, err)
			return
		}
		info := model.StockInfo{
			ID:        id,
			Name:      name,
			FirstDate: start.Date,
			LastDate:  last.Date,
			Unable:    !last.Date.Equal(lastOpen.Date),
		}
		if last.Date.Equal(lastOpen.Date) {
			svc.l.Debugf("%s - %s Enable", id, name)
		}
		if err := db.Model(&info).Updates(info).Error; err != nil {
			svc.l.Errorf("update stock list failed, %+v", err)
			return
		}
	}
}
