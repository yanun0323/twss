package service

import (
	"errors"
	"stocker/internal/util"
	"time"
)

type CheckRawOption struct {
	CrawlOpt CrawlDateOption

	Name         string
	GetBeginDate func(Service) (time.Time, error)
	GetRaw       func(Service, time.Time) (interface{}, error)
}

// CheckRawTrade 檢查 RawTrade 資料是否有缺漏
var CheckRawTrade = CheckRawOption{
	CrawlOpt: CrawlTrade,

	Name: "check_trade",
	GetBeginDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawTradeDate(svc.Ctx, true)
	},
	GetRaw: func(svc Service, date time.Time) (interface{}, error) {
		return svc.Repo.GetRawTrade(svc.Ctx, date)
	},
}

// CheckRawFinance 檢查 RawFinance 資料是否有缺漏
var CheckRawFinance = CheckRawOption{
	CrawlOpt: CrawlFinance,

	Name: "check_finance",
	GetBeginDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawFinanceDate(svc.Ctx, true)
	},
	GetRaw: func(svc Service, date time.Time) (interface{}, error) {
		return svc.Repo.GetRawFinance(svc.Ctx, date)
	},
}

// CheckRaw 檢查 Raw 資料是否有缺漏
func (svc Service) CheckRaw(repair bool, opt CheckRawOption) {
	svc.Log = svc.Log.WithField("service", opt.Name)
	date, err := opt.GetBeginDate(svc)
	if err != nil {
		svc.Log.Errorf("get default start date , %+v", err)
		return
	}
	count := 0

	date = date.Add(24 * time.Hour)
	now := time.Now().Local().Add(opt.CrawlOpt.TimeOffset)
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		count++
		logDate := util.LogDate(date)
		_, err := opt.GetRaw(svc, date)
		if err != nil {
			svc.Log.Errorf("%s, get raw, err: %+v", logDate, err)
			return
		}
		if errors.Is(svc.Repo.ErrNotFound(), err) {
			svc.Log.Errorf("%s, found missing raw", logDate)
			if !repair {
				continue
			}
			svc.Log.Infof("repairing %s raw", logDate)
			err := svc.crawlRaw(svc.Ctx, date, opt.CrawlOpt)
			if err != nil {
				svc.Log.Errorf("%s, crawl raw, err: %+v", logDate, err)
			}
		}
	}
	svc.Log.Infof("checked raw count: %d", count)
	svc.Log.Info("check raw done!")
}
