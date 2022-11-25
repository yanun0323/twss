package service

import (
	"fmt"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"
)

const (
	_API_LIMIT_INTERVAL_TIME = 3 * time.Second
	_REQUEST_RETRY_TIMES     = 3
)

func (svc Service) CrawlDailyRawData() {
	svc.l.Info("start daily raw data crawl")
	last, err := svc.repo.GetLastDailyRawDate()
	if err != nil {
		svc.l.Errorf("get last daily raw date failed, %+v", err)
		return
	}

	date := last.Add(24 * time.Hour)
	now := time.Now().Local().Add(-18 * time.Hour) /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */

	svc.l.Debug("start crawl date ", util.LogDate(date))
	svc.l.Debug("start crawl now  ", util.LogDate(now))
	for ; date.Before(now); util.NextDay(&date) {
		for r := _REQUEST_RETRY_TIMES; r > 0; r-- {
			err := svc.crawl(date)
			if err != nil {
				break
			}
			if r > 1 {
				svc.l.Warnf("crawl failed, retry in 3 second, remain %d times, %+v", r, err)
			}
		}
		svc.l.Errorf("crawl date %s failed, stop crawling", util.LogDate(date))
	}
	svc.l.Info("all daily raw data is update to date")
}

func (svc Service) crawl(date time.Time) error {
	defer time.Sleep(_API_LIMIT_INTERVAL_TIME)
	logDate := util.LogDate(date)
	svc.l.Infof("--- start crawl %s ---", logDate)
	
	url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=%s&type=ALLBUT0999", util.FormatToUrlDate(date))
	body, err := util.GetRequest(url)
	if err != nil {
		return err
	}

	raw := model.DailyRaw{
		Date: date,
		Body: body,
	}

	if err := svc.repo.InsertDailyRaw(raw); err != nil {
		return err
	}
	svc.l.Infof("crawl success %s, data size: %d", logDate, len(body))
	return nil
}
