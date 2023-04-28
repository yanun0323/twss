package service

import (
	"errors"
	"stocker/internal/util"
	"time"
)

func (svc Service) CheckDailyRaw() {
	svc.l.Info("start checking daily raw ...")
	date, err := svc.repo.GetDefaultStartDate()
	if err != nil {
		svc.l.Errorf("get default start date , %+v", err)
		return
	}
	count := 0
	now := date.UTC()
	for ; date.Before(now); date = date.UTC().Add(24 * time.Hour) {
		count++
		_, err := svc.repo.GetRawTrade(date)
		svc.l.Infof("Date: %s", date.Format("2006-01-02 15:04:05 Z07:00"))
		if errors.Is(svc.repo.ErrRecordNotFound(), err) {
			svc.l.Errorf("%s, found missing daily raw", util.LogDate(date))
		}
	}
	svc.l.Infof("checked data count: %d", count)
	svc.l.Info("check daily raw done!")
}

func (svc Service) CheckConverter() {
	svc.l.Info("start checking daily raw data converter ...")
	date, err := svc.repo.GetDefaultStartDate()
	if err != nil {
		svc.l.Errorf("get default start date , %+v", err)
		return
	}
	count := 0
	now := time.Now().Local().Add(-18 * time.Hour)
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		count++
		if err := svc.repo.CheckOpen(date); errors.Is(svc.repo.ErrRecordNotFound(), err) && err != nil {
			svc.l.Errorf("%s, found missing open date", util.LogDate(date))
			continue
		}
		if err := svc.repo.CheckStock(date); errors.Is(svc.repo.ErrRecordNotFound(), err) && err != nil {
			svc.l.Errorf("%s, found missing stock date", util.LogDate(date))
			continue
		}
	}
	svc.l.Infof("checked data count: %d", count)
	svc.l.Info("check daily raw data converter done!")
}
