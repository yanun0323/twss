package service

import (
	"context"
	"net/http"
	"net/url"
	"stocker/internal/model"
	"strings"
	"time"
)

func (su *ServiceTestSuite) TestCrawl() {
	ctx := context.WithValue(su.ctx, "debug", true)
	date, err := time.ParseInLocation("20060102", "20230501", time.Local)
	su.Require().Nil(err)
	su.Assert().Nil(su.svc.crawlRaw(ctx, date, CrawlProfit))
}

var CrawlProfit = CrawlDateOption{
	Name:            "crawl_profit",
	Method:          http.MethodPost,
	UrlFormat:       "https://mops.twse.com.tw/mops/web/ajax_t163sb06",
	UrlFormatArgsFn: nil,
	TimeOffset:      -18 * time.Hour, /* turn every 18:00 into 00:00 to crawl data after 18:00 every day */
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
		return nil
	},
	RequestFn: func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Origin", "https://mops.twse.com.tw")
		req.Header.Set("Host", "mops.twse.com.tw")
		req.Header.Set("Cookie", "jcsession=jHttpSession@714fb264")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6.1 Safari/605.1.15")
		req.Header.Set("Accept-Language", "zh-TW,zh-Hant;q=0.9")
		req.Header.Set("Connection", "keep-alive")
	},
	RequestBody: strings.NewReader(url.Values{
		"encodeURIComponent": []string{"1"},
		"step":               []string{"1"},
		"firstin":            []string{"1"},
		"off":                []string{"1"},
		"isQuery":            []string{"Y"},
		"TYPEK":              []string{"sii"},
		"year":               []string{"112"},
		"season":             []string{"1"},
	}.Encode()),
}
