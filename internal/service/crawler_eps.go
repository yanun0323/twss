package service

import (
	"fmt"
	"stocker/internal/model"
	"stocker/internal/util"
	"time"

	"github.com/pkg/errors"
)

func (svc Service) CrawlRawEps() {
	svc.l.Info("start eps raw data crawl ...")
	last, err := svc.repo.GetLastRawEpsDate(svc.ctx)
	if err != nil {
		svc.l.Errorf("get last eps raw date , %+v", err)
		return
	}

	date := last.Add(24 * time.Hour)
	now := time.Now().Local().Add(-18 * time.Hour) /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */

	svc.l.Debugf("start crawl eps date: %s", util.LogDate(date))
	svc.l.Debugf("start crawl eps now: %s", util.LogDate(now))
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		for r := _REQUEST_RETRY_TIMES; r >= 0; r-- {
			err := svc.crawlRawEps(date)
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
	svc.l.Info("all eps raw data is update to date!")
}

func (svc Service) crawlRawEps(date time.Time) error {
	defer time.Sleep(_API_LIMIT_INTERVAL_TIME)
	logDate := util.LogDate(date)
	svc.l.Infof("--- start crawl %s ---", logDate)

	url := fmt.Sprintf("https://www.twse.com.tw/rwd/zh/afterTrading/BWIBBU_d?date=%s&selectType=ALL&response=json", util.FormatToUrlDate(date))
	body, err := util.GetRequest(url)
	if err != nil {
		return errors.Errorf("get request, err: %+v", err)
	}

	raw := model.RawEps{
		Date: date,
		Body: body,
	}

	if err := svc.repo.InsertRawEps(svc.ctx, raw); err != nil {
		return errors.Errorf("insert raw eps, err: %+v", err)
	}
	svc.l.Infof("crawl success %s, data size: %d", logDate, len(body))
	return nil
}
