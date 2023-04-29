package service

import (
	"stocker/internal/model"
	"time"
)

func (svc Service) Debug() {
	svc.Log.Info("start debugging ...")

	if false { /* to avoid unused warning */
		svc.debugTradeRawData("20221125")
	}
}

func (svc Service) debugTradeRawData(dateStr string) {
	date, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	raw, err := svc.Repo.GetRawTrade(svc.Ctx, date)
	if err != nil {
		svc.Log.Errorf("get trade raw , %+v", err)
		return
	}
	data, err := raw.GetData()
	if err != nil {
		svc.Log.Errorf("get data , %+v", err)
		return
	}
	d, ok := data.(model.RawTradeData)
	if !ok {
		svc.Log.Errorf("type assertion , %+v", err)
		return
	}

	svc.Log.Debugf("%+v", d.Date)
	svc.Log.Debugf("%+v", d.Fields8)
	svc.Log.Debugf("%+v", d.TradeData()[0])
	svc.Log.Debugf("%+v", d.Parse()[0])
}
