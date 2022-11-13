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

func (svc Service) CrawlRawDaily() {
	last, err := svc.repo.GetLastRaw()
	if err != nil {
		svc.l.Errorf("failed to get last raw, %+v", err)
		return
	}
	date := last.Date.Add(24 * time.Hour)

	now := time.Now().Local().Add(-18 * time.Hour) /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */
	if date.Before(now) {
		svc.l.Info("start crawl daily raw data, last raw date: ", util.LogDate(last.Date))
	}
	svc.l.Debug("now ", now)
	for ; date.Before(now); util.NextDay(&date) {
		for r := _REQUEST_RETRY_TIMES; r > 0; r-- {
			err := svc.crawl(date)
			if err == nil {
				break
			}
			svc.l.Errorf("crawl failed, retry in 3 second, %+v", err)
		}
	}
	svc.l.Info("all raw data is update to date")
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

	raw := model.Raw{
		Date: date,
		Body: string(body),
	}

	if err := svc.repo.InsertRaw(raw); err != nil {
		return err
	}
	svc.l.Infof("crawl success %s, data size: %d", logDate, len(body))
	return nil
}
