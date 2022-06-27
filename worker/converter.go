package worker

import (
	"encoding/json"
	"log"
	"main/domain"
	"main/model"
	"main/setting"
	"main/util"
	"sync"
	"time"
)

const _MAX_GOROUTINE = 100

type Converter struct {
	repo         domain.IRepository
	date         time.Time
	stockHash    chan map[string]string
	parseChannel chan int8
	checkMode    bool
}

func NewConverter(repo domain.IRepository, checkMode bool) *Converter {
	return &Converter{
		repo:         repo,
		stockHash:    make(chan map[string]string, 1),
		parseChannel: make(chan int8, _MAX_GOROUTINE),
		checkMode:    checkMode,
	}
}
func (c *Converter) InitMigrate() {
	err := c.repo.AutoMigrate(&model.OpenDays{}, &model.StockList{})
	if err != nil {
		log.Println(err)
	}
}

func (c *Converter) Run() {
	date := c.repo.GetConvertibleDate(c.checkMode)

	c.date = date
	c.stockHash <- c.repo.GetStockHash()
	for {
		if c.date.Add(setting.Time_Offset).After(time.Now().Local()) {
			log.Printf("over time: %s", util.LogDate(c.date.Add(setting.Time_Offset)))
			return
		}
		log.Println(util.LogDate(c.date))
		if c.checkMode {
			log.Print("- checkMode")
		}

		body, err := c.repo.GetRaw(c.date)
		if err != nil {
			log.Println(err)
			break
		}
		var raw model.RawJson
		_ = json.Unmarshal(body, &raw)
		if raw.State != "OK" {
			c.repo.Insert(model.OpenDays{Date: c.date, State: false})
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
		err = c.repo.Insert(model.OpenDays{Date: c.date, State: true})
		if err != nil {
			log.Printf("failed to insert open day %s %s", util.LogDate(c.date), err)
		}
		c.date = util.NextDate(c.date)
		log.Println("complete")
		// return
	}
}

func (c *Converter) ParseService(wg *sync.WaitGroup, d []string) {
	defer wg.Done()
	c.parseChannel <- 0
	id := d[0]
	name := d[1]
	table := util.StockTable(id)

	hash := <-c.stockHash
	_, exist := hash[id]
	if !exist {
		hash[id] = name
		c.repo.Migrate(table, &model.Deal{})
		c.repo.Insert(
			model.StockList{
				StockID:   id,
				StockName: name,
			},
		)
	}
	c.stockHash <- hash
	deal := model.NewDealFromDataString(d, c.date)
	err := c.repo.InsertWithTableName(table, &deal)
	if err != nil {
		log.Printf("Failed to insert %s %s %s", id, d[1], err)
	}

	<-c.parseChannel
}
