package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/config"
	"main/domain"
	"main/model"
	"main/util"
	"net/http"
	"time"
)

const waitSecond = 3 * time.Second

type Crawler struct {
	repo      domain.IRepository
	date      time.Time
	checkMode bool
}

func NewCrawler(repo domain.IRepository, checkMode bool) *Crawler {
	d := repo.GetCrawlableDate(checkMode)
	return &Crawler{repo: repo, date: d, checkMode: checkMode}
}

func (c *Crawler) InitMigrate() {
	c.repo.AutoMigrate(&model.Raw{})
}

func (c *Crawler) Run() {
	for {
		if c.date.Add(config.TimeOffset).After(time.Now().Local()) {
			log.Printf("over time: %s", util.LogDate(c.date.Add(config.TimeOffset)))
			break
		}
		log.Println(util.LogDate(c.date))
		if c.checkMode {
			log.Print("- checkMode")
		}

		body := c.crawlPrice(c.date)

		if len(body) == 0 {
			c.date = util.NextDate(c.date)
			time.Sleep(waitSecond)
			continue
		}

		err := c.repo.Insert(model.Raw{Date: c.date, Body: body})
		if err != nil {
			log.Println(err)
		}

		c.date = util.NextDate(c.date)
		time.Sleep(waitSecond)
	}
	log.Println("Crawl complete")
}

func (c *Crawler) crawlPrice(target time.Time) []byte {
	date := util.FormatDate(target)
	url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=%s&type=ALLBUT0999", date)
	log.Println(url)

	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("get content failed status code is %d. \n", response.StatusCode)
		return []byte{}
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("HTML body read error")
		return []byte{}
	}
	log.Printf("Get HTML. %s", util.LogDate(target))

	return bytes
}
