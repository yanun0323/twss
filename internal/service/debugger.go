package service

import (
	"stocker/internal/model"
	"stocker/internal/util"
	"time"
)

func (svc Service) Debug() {
	svc.l.Info("start debugging ...")
}

func (svc Service) debugConvertRawToDailyStock() {
	begin, _ := svc.repo.GetDefaultStartDate()
	raws, err := svc.repo.ListDailyRaws(begin, time.Now())
	if err != nil {
		svc.l.Errorf("get daily raw failed, %+v", err)
		return
	}

	for _, raw := range raws {
		rawData, err := raw.GetData()
		if err != nil {
			svc.l.Errorf("%s, get data failed, %+v", util.LogDate(raw.Date), err)
			return
		}

		data, _, err := rawData.ParseStock()
		if err != nil {
			svc.l.Errorf("%s, parse stock failed, %+v", util.LogDate(raw.Date), err)
			return
		}

		if err := svc.repo.InsertDailyStock(data); err != nil {
			svc.l.Errorf("%s, create daily stock failed, %+v", util.LogDate(raw.Date), err)
			return
		}
		svc.l.Infof("%s, convert succeed", util.LogDate(raw.Date))
	}
	svc.l.Info("all raw converted!")
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
	_, data, err := d.ParseStock()
	if err != nil {
		svc.l.Errorf("parse stock failed, %+v", err)
		return
	}

	svc.l.Debugf("%+v", d.RawDate)
	svc.l.Debugf("%+v", d.Fields8)
	svc.l.Debugf("%+v", d.Data()[0])
	svc.l.Debugf("%+v", data[0])
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
		table := model.DailyStockData{ID: id}.GetTableName()
		start, last := model.DailyStockData{}, model.DailyStockData{}
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
