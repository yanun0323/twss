package service

import (
	"stocker/internal/model"
	"stocker/internal/util"
	"sync"
	"time"
)

func (svc Service) ConvertDailyRawData() {
	svc.l.Info("start daily raw data convert ...")
	stockMap, err := svc.repo.GetStockMap()
	if err != nil {
		svc.l.Errorf("get stock map , %+v", err)
		return
	}

	date, err := svc.repo.GetLastOpenDate()
	if err != nil {
		svc.l.Errorf("get last open date , %+v", err)
		return
	}
	date = date.Add(24 * time.Hour)

	raws, err := svc.repo.ListRawTrades(date, time.Now())
	if err != nil {
		svc.l.Errorf("list daily raws , %+v", err)
		return
	}

	inserterWP := util.NewWorkerPool("InsertDailyStock", 150)
	inserterWP.Run()
	stockChan := make(chan model.Stock, 150)
	closeLooperChan := make(chan struct{}, 1)
	var looperWG sync.WaitGroup
	looperWG.Add(1)
	go func() {
		defer looperWG.Done()
		for {
			select {
			case stock := <-stockChan:
				inserterWP.Push(func() {
					err := svc.repo.InsertStock(stock)
					if err != nil {
						svc.l.Errorf("%s, insert stock , %+v", util.LogDate(stock.Date), err)
					}
				})
			case <-closeLooperChan:
				svc.l.Debug("insert looper stopped")
				return
			}
		}
	}()
	svc.l.Debug("raws count:", len(raws))
	for _, raw := range raws {
		svc.convert(stockMap, raw, stockChan)
	}

	svc.l.Info("starting shutdown inserter worker pool ...")
	if err := inserterWP.Shutdown(30 * time.Second); err != nil {
		svc.l.Errorf("shutdown inserter worker pool with error, %+v", err)
	}
	closeLooperChan <- struct{}{}
	looperWG.Wait()

	svc.l.Info("all daily raw data converted!")
}

func (svc Service) convert(stockMap model.StockMap, raw model.RawTrade, stockChan chan model.Stock) {
	logDate := util.LogDate(raw.Date)
	data, err := raw.GetData()
	if err != nil {
		svc.l.Errorf("%s, unmarshal daily raw data , %+v", logDate, err)
		return
	}

	if data.Stat != "OK" || len(data.StockData()) == 0 {
		_ = svc.repo.InsertOpen(model.Open{
			Date:   raw.Date,
			IsOpen: false,
		})
		return
	}

	stocks := data.ParseStock()
	for _, stock := range stocks {
		if _, exist := stockMap[stock.ID]; !exist {
			stockMap[stock.ID] = stock.Name
			_ = svc.repo.InsertStockList(model.StockListUnit{
				ID:   stock.ID,
				Name: stock.Name,
			})
		}
		stockChan <- stock
	}
	_ = svc.repo.InsertOpen(model.Open{
		Date:   raw.Date,
		IsOpen: true,
	})
	svc.l.Infof("%s, convert succeed", logDate)
}
