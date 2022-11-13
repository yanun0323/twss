package service

import (
	"encoding/json"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"
)

func (svc Service) ConvertDailyRawData() {
	raws, err := svc.repo.ListAllDailyRaws()
	if err != nil {
		svc.l.Errorf("list all daily raws failed, %+v", err)
		return
	}

	svc.l.Infof("start convert daily raw data, daily raws amount: %d", len(raws))

	wp := util.NewWorkerPool("ConvertDailyRawData", 15)
	wp.Run()
	for _, raw := range raws {
		func(r model.DailyRaw) {
			wp.Push(func() {
				svc.convert(r)
			})
		}(raw)
	}

	if err := wp.Shutdown(30 * time.Second); err != nil {
		svc.l.Errorf("shutdown worker pool with error, %+v", err)
	}

	svc.l.Info("all daily raw data converted")
}

func (svc Service) convert(raw model.DailyRaw) {
	logDate := util.LogDate(raw.Date)
	data := &model.DailyRawData{}
	if err := json.Unmarshal([]byte(raw.Body), data); err != nil {
		svc.l.Errorf("%s, unmarshal daily raw data failed, %+v", logDate, err)
		return
	}

	if data.Stat != "OK" || len(data.Data) == 0 {
		return
	}
	// TODO: Update `Open` and `Stock List`
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
