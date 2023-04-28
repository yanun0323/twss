package service

import (
	"stocker/internal/util"
	"time"
)

// CheckTradeRaw 檢查 TradeRaw 資料是否有缺漏
func (svc Service) CheckTradeRaw(repair bool) {
	svc.l.Info("start checking trade raw ...")
	date, err := svc.repo.GetDefaultStartDate()
	if err != nil {
		svc.l.Errorf("get default start date , %+v", err)
		return
	}
	count := 0
	now := date.UTC()
	for ; date.Before(now); date = date.UTC().Add(24 * time.Hour) {
		count++
		logDate := util.LogDate(date)
		_, err := svc.repo.GetRawTrade(svc.ctx, date)
		svc.l.Infof("Date: %s", date.Format("2006-01-02 15:04:05 Z07:00"))
		if svc.repo.IsErrRecordNotFound(err) {
			svc.l.Errorf("%s, found missing trade raw", logDate)
			if !repair {
				continue
			}
			svc.l.Infof("repairing %s trade raw", logDate)
			err := svc.crawlRawTrade(date)
			if err != nil {
				svc.l.Errorf("%s, crawl raw trade, err: %+v", logDate, err)
			}
		}
	}
	svc.l.Infof("checked data count: %d", count)
	svc.l.Info("check trade raw done!")
}

// CheckTradeRaw 檢查 Trade 資料是否有缺漏
func (svc Service) CheckConverter(repair bool) {
	svc.l.Info("start checking trade raw data converter ...")
	date, err := svc.repo.GetDefaultStartDate()
	if err != nil {
		svc.l.Errorf("get default start date , %+v", err)
		return
	}
	count := 0
	now := time.Now().Local().Add(-18 * time.Hour)
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		count++
		logDate := util.LogDate(date)
		if err := svc.repo.CheckTradeDate(svc.ctx, date); svc.repo.IsErrRecordNotFound(err) && err != nil {
			svc.l.Errorf("%s, found missing trade date", logDate)
			if !repair {
				continue
			}
			svc.l.Infof("repairing %s trade date", logDate)
			err := svc.convertRawTrade(date)
			if err != nil {
				svc.l.Errorf("%s, repair trade, err: %+v", logDate, err)
				continue
			}
		}

		if err := svc.repo.CheckTrade(svc.ctx, date); svc.repo.IsErrRecordNotFound(err) && err != nil {
			svc.l.Errorf("%s, found missing stock date, unrepairable", logDate)
			continue
		}
	}
	svc.l.Infof("checked data count: %d", count)
	svc.l.Info("check trade raw data converter done!")
}
