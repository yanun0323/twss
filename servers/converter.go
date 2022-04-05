package servers

import (
	"encoding/json"
	"log"
	"main/domain"
	"main/model"
	"main/util"
	"sync"
	"time"
)

const maxGoroutine = 100

type Converter struct {
	repo         domain.IRepository
	date         time.Time
	stockHash    chan map[string]string
	maxGoroutine chan int8
}

func NewConverter(repo domain.IRepository) *Converter {
	return &Converter{repo: repo}
}
func (c *Converter) InitMigrate() {
	err := c.repo.AutoMigrate(&model.OpenDays{}, &model.StockList{})
	if err != nil {
		log.Println(err)
	}
}

func (c *Converter) Run() {

	c.stockHash = make(chan map[string]string, 1)
	c.maxGoroutine = make(chan int8, maxGoroutine)
	date, ok := c.repo.GetConvertableDate()
	if !ok {
		log.Println("Failed to get convertable date")
		return
	}

	c.date = date
	c.stockHash <- c.repo.GetStockHash()
	for {
		log.Println(util.LogDate(c.date))

		body, err := c.repo.GetRaw(c.date)
		if err != nil {
			log.Println(err)
			break
		}
		var raw model.RawJson
		_ = json.Unmarshal(body, &raw)
		if raw.State != "OK" {
			c.repo.Insert(&model.OpenDays{Date: c.date, State: false})
			c.date = util.NextDate(c.date)
			continue
		}

		// data8: before 2011/7/31
		// data9: '2006/09/29' and after 2011/7/31
		data := raw.Data8
		if len(data) < len(raw.Data9) {
			data = raw.Data9
		}
		length := len(data)
		log.Printf("%d stocks start converting...", length)
		var wg sync.WaitGroup
		wg.Add(length)
		for _, d := range data {
			go c.ParseService(&wg, d)
		}
		wg.Wait()
		err = c.repo.Insert(&model.OpenDays{Date: c.date, State: true})
		if err != nil {
			log.Printf("Failed to insert open day %s %s", util.LogDate(c.date), err)
		}
		c.date = util.NextDate(c.date)
		log.Println("Complete")
		// return
	}
}

func (c *Converter) ParseService(wg *sync.WaitGroup, d []string) {
	defer wg.Done()
	c.maxGoroutine <- 0
	id := d[0]
	name := d[1]
	table := util.StockTable(id)

	hash := <-c.stockHash
	_, exist := hash[id]
	if !exist {
		hash[id] = name
		c.repo.Migrate(table, &model.Deal{})
		c.repo.Insert(
			&model.StockList{
				StockID:   id,
				StockName: name,
			},
		)
	}
	c.stockHash <- hash
	deal := &model.Deal{
		Date:        c.date,
		Volume:      d[2],
		VolumeMoney: d[3],
		Start:       d[5],
		Max:         d[6],
		Min:         d[7],
		End:         d[8],
		Grade:       d[9],
		Spread:      d[10],
		Per:         d[15],
	}

	err := c.repo.Create(table, &deal)
	if err != nil {
		log.Printf("Failed to insert %s %s %s", id, d[1], err)
	}

	<-c.maxGoroutine
}
