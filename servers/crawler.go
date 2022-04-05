package servers

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/domain"
	"main/model"
	"main/util"
	"net/http"
	"time"
)

type Crawler struct {
	repo domain.IRepository
	date time.Time
}

func NewCrawler(repo domain.IRepository) *Crawler {
	d := repo.GetCrawlableDate()
	return &Crawler{repo: repo, date: d}
}

func (c *Crawler) InitMigrate() {
	c.repo.AutoMigrate(&model.Raw{})
}

func (c *Crawler) Run() {
	for {
		if c.date.After(time.Now().Local()) {
			break
		}
		log.Println(util.LogDate(c.date))
		body := c.crawlPrice()

		if len(body) == 0 {
			c.date = util.NextDate(c.date)
			time.Sleep(2000 * time.Millisecond)
			continue
		}

		err := c.repo.Insert(&model.Raw{Date: c.date, Body: body})
		if err != nil {
			log.Println(err)
		}

		c.date = util.NextDate(c.date)
		time.Sleep(2000 * time.Millisecond)
	}
	log.Println("Crawl complete")
}

func (c *Crawler) crawlPrice() []byte {
	date := util.FormatDate(c.date)
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
	log.Printf("Get HTML. %s", util.LogDate(c.date))

	return bytes
}
