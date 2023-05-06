package service

import (
	"errors"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"
)

type CheckDataOption struct {
	ConvertOpt ConvertRawOption

	Name         string
	GetBeginDate func(Service) (time.Time, error)
	GetDataDate  func(Service, time.Time) (model.DataDate, error)
	IsDataExist  func(Service, time.Time) (bool, error)
}

// CheckTrade 檢查 Trade 資料是否有缺漏
var CheckTrade = CheckDataOption{
	ConvertOpt: ConvertRawTrade,

	Name: "check_trade",
	GetBeginDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawTradeDate(svc.Ctx, true)
	},
	GetDataDate: func(svc Service, date time.Time) (model.DataDate, error) {
		return svc.Repo.GetTradeDate(svc.Ctx, date)
	},
	IsDataExist: func(svc Service, date time.Time) (bool, error) {
		return svc.Repo.IsTradeExist(svc.Ctx, date)
	},
}

// CheckFinance 檢查 Finance 資料是否有缺漏
var CheckFinance = CheckDataOption{
	ConvertOpt: ConvertRawFinance,

	Name: "check_finance",
	GetBeginDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawFinanceDate(svc.Ctx, true)
	},
	GetDataDate: func(svc Service, date time.Time) (model.DataDate, error) {
		return svc.Repo.GetFinanceDate(svc.Ctx, date)
	},
	IsDataExist: func(svc Service, date time.Time) (bool, error) {
		return svc.Repo.IsFinanceExist(svc.Ctx, date)
	},
}

// CheckData 檢查資料是否有缺漏
func (svc Service) CheckData(repair bool, opt CheckDataOption) {
	svc.Log = svc.Log.WithField("service", opt.Name)
	svc.Log.Debug("start checking converted data...")
	date, err := opt.GetBeginDate(svc)
	if err != nil {
		svc.Log.Errorf("get begin date, %+v", err)
		return
	}
	count := 0

	date = date.Add(24 * time.Hour)
	now := time.Now().Local().Add(opt.ConvertOpt.TimeOffset)
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		count++
		logDate := util.LogDate(date)
		dataDate, err := opt.GetDataDate(svc, date)

		if errors.Is(svc.Repo.ErrNotFound(), err) {
			svc.Log.Errorf("%s, found missing data date", logDate)
			if !repair {
				continue
			}
			svc.Log.Infof("repairing %s data date", logDate)
			err := svc.convertRaw(date, opt.ConvertOpt)
			if err != nil {
				svc.Log.Errorf("%s, repair trade, err: %+v", logDate, err)
				continue
			}
			continue
		}

		if err != nil {
			svc.Log.Errorf("%s, get trade date, err: %+v", logDate, err)
			return
		}

		if !dataDate.IsOpen() {
			continue
		}

		exist, err := opt.IsDataExist(svc, date)
		if err != nil {
			svc.Log.Errorf("%s, check trade, err: %+v", logDate, err)
			return
		}

		if !exist {
			svc.Log.Errorf("%s, found missing stock date, unrepairable", logDate)
			continue
		}
	}

	if repair {
		if err := svc.storeStockMap(); err != nil {
			svc.Log.Errorf("store stock map, err: %+v", err)
		}
	}

	svc.Log.Infof("checked data count: %d", count)
	svc.Log.Info("check trade raw data converter done!")
}
