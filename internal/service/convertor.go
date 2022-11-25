package service

import (
	"encoding/json"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"
)

func (svc Service) ConvertDailyRawData() {
	svc.l.Info("start daily raw data convert")
	stockMap, err := svc.repo.GetStockMap()
	if err != nil {
		svc.l.Errorf("get stock map failed, %+v", err)
		return
	}

	date, err := svc.repo.GetLastOpenDate()
	if err != nil {
		svc.l.Errorf("get stock map failed, %+v", err)
		return
	}
	date = date.Add(24 * time.Hour)

	stockMapChan := make(chan model.StockMap, 1)
	stockMapChan <- stockMap

	raws, err := svc.repo.ListDailyRaws(date, time.Now())
	if err != nil {
		svc.l.Errorf("list daily raws failed, %+v", err)
		return
	}
	wp := util.NewWorkerPool("ConvertDailyRawData", 15)
	wp.Run()
	// TODO: push insert stock to worker pool instead push convert function
	for _, raw := range raws {
		func(r model.DailyRaw) {
			wp.Push(func() {
				svc.convert(stockMapChan, r)
			})
		}(raw)
	}

	if err := wp.Shutdown(30 * time.Second); err != nil {
		svc.l.Errorf("shutdown worker pool with error, %+v", err)
	}

	svc.l.Info("all daily raw data converted")
}

func (svc Service) convert(stockMapChan chan model.StockMap, raw model.DailyRaw) {
	logDate := util.LogDate(raw.Date)
	data := &model.DailyRawData{}
	if err := json.Unmarshal([]byte(raw.Body), data); err != nil {
		svc.l.Errorf("%s, unmarshal daily raw data failed, %+v", logDate, err)
		return
	}

	if data.Stat != "OK" || len(data.Data) == 0 {
		return
	}

	stocks := data.ParseStock(raw.Date)
	for _, stock := range stocks {
		if stock.ID != "2330" {
			continue
		}
		err := svc.repo.InsertDailyStock(stock)
		if err != nil {
			svc.l.Errorf("%s, insert stock failed, %+v", logDate, err)
		}
		return
	}
}
