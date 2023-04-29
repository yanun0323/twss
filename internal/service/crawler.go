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

type CrawlDateOption struct {
	Name           string
	UrlFormat      string
	TimeOffset     time.Duration
	GetLastRawDate func(Service) (time.Time, error)
	CreateRaw      func(time.Time, []byte) interface{}
	InsertRaw      func(Service, interface{}) error
}

// CrawlTrade 爬蟲每日盤後交易資料
var CrawlTrade = CrawlDateOption{
	Name:       "crawl_trade",
	UrlFormat:  "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=%s&type=ALLBUT0999",
	TimeOffset: -18 * time.Hour, /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */
	GetLastRawDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawTradeDate(svc.Ctx, false)
	},
	CreateRaw: func(date time.Time, body []byte) interface{} {
		return model.RawTrade{
			Date: date,
			Body: body,
		}
	},
	InsertRaw: func(svc Service, obj interface{}) error {
		raw, ok := obj.(model.RawTrade)
		if !ok {
			return errors.New("invalid raw trade type")
		}
		return svc.Repo.InsertRawTrade(svc.Ctx, raw)
	},
}

// CrawlEps 爬蟲每日EPS
var CrawlEps = CrawlDateOption{
	Name:       "crawl_eps",
	UrlFormat:  "https://www.twse.com.tw/rwd/zh/afterTrading/BWIBBU_d?date=%s&selectType=ALL&response=json",
	TimeOffset: -18 * time.Hour, /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */
	GetLastRawDate: func(svc Service) (time.Time, error) {
		return svc.Repo.GetRawEpsDate(svc.Ctx, false)
	},
	CreateRaw: func(date time.Time, body []byte) interface{} {
		return model.RawEps{
			Date: date,
			Body: body,
		}
	},
	InsertRaw: func(svc Service, obj interface{}) error {
		raw, ok := obj.(model.RawEps)
		if !ok {
			return errors.New("invalid raw trade type")
		}
		return svc.Repo.InsertRawEps(svc.Ctx, raw)
	},
}

func (svc Service) CrawlRaw(opt CrawlDateOption) {
	svc.Log = svc.Log.WithField("service", opt.Name)
	svc.Log.Info("start raw crawl ...")
	last, err := opt.GetLastRawDate(svc)
	if err != nil {
		svc.Log.Errorf("get last raw date , %+v", err)
		return
	}

	date := last.Add(24 * time.Hour)
	now := time.Now().Local().Add(opt.TimeOffset) /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */

	svc.Log.Debugf("start crawl date: %s", util.LogDate(date))
	svc.Log.Debugf("start crawl now: %s", util.LogDate(now))
	for ; date.Before(now); date = date.Add(24 * time.Hour) {
		for r := _REQUEST_RETRY_TIMES; r >= 0; r-- {
			err := svc.crawlRaw(date, opt)
			if err == nil {
				break
			}
			if r == 0 {
				svc.Log.Errorf("crawl date %s failed, stop crawling", util.LogDate(date))
				return
			}
			svc.Log.Warnf("crawl failed, retry in 3 second, remain %d times, %+v", r, err)
		}
	}
	svc.Log.Info("all raw is update to date!")
}

func (svc Service) crawlRaw(date time.Time, opt CrawlDateOption) error {
	defer time.Sleep(_API_LIMIT_INTERVAL_TIME)
	logDate := util.LogDate(date)
	svc.Log.Infof("--- start crawl %s ---", logDate)
	url := fmt.Sprintf(opt.UrlFormat, util.FormatToUrlDate(date))
	body, err := util.GetRequest(url)
	if err != nil {
		return errors.Errorf("get request, err: %+v", err)
	}

	raw := opt.CreateRaw(date, body)

	if err := opt.InsertRaw(svc, raw); err != nil {
		return errors.Errorf("insert raw, err: %+v", err)
	}
	svc.Log.Infof("crawl success %s, data size: %d", logDate, len(body))
	return nil
}
