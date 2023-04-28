package service

import (
	"fmt"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"

	"github.com/pkg/errors"
)

const (
	_API_LIMIT_INTERVAL_TIME = 3 * time.Second
	_REQUEST_RETRY_TIMES     = 3
)

func (svc Service) CrawlRawTrade() {
	svc.l.Info("start trade raw data crawl ...")
	last, err := svc.repo.GetLastRawTradeDate(svc.ctx)
	if err != nil {
		svc.l.Errorf("get last trade raw date , %+v", err)
		return
	}

	date := last.Add(24 * time.Hour)
	now := time.Now().Local().Add(-18 * time.Hour) /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */

	svc.l.Debugf("start crawl trade date: %s", util.LogDate(date))
	svc.l.Debugf("start crawl trade now: %s", util.LogDate(now))
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		for r := _REQUEST_RETRY_TIMES; r >= 0; r-- {
			err := svc.crawlRawTrade(date)
			if err == nil {
				break
			}
			if r == 0 {
				svc.l.Errorf("crawl date %s failed, stop crawling", util.LogDate(date))
				return
			}
			svc.l.Warnf("crawl failed, retry in 3 second, remain %d times, %+v", r, err)
		}
	}
	svc.l.Info("all trade raw data is update to date!")
}

func (svc Service) crawlRawTrade(date time.Time) error {
	defer time.Sleep(_API_LIMIT_INTERVAL_TIME)
	logDate := util.LogDate(date)
	svc.l.Infof("--- start crawl %s ---", logDate)
	//https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=20230428&type=ALLBUT0999
	url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=%s&type=ALLBUT0999", util.FormatToUrlDate(date))
	body, err := util.GetRequest(url)
	if err != nil {
		return errors.Errorf("get request, err: %+v", err)
	}

	raw := model.RawTrade{
		Date: date,
		Body: body,
	}

	if err := svc.repo.InsertRawTrade(svc.ctx, raw); err != nil {
		return errors.Errorf("insert raw trade, err: %+v", err)
	}
	svc.l.Infof("crawl success %s, data size: %d", logDate, len(body))
	return nil
}
