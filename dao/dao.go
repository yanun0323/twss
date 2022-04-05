package dao

import (
	"encoding/json"
	"fmt"
	"log"
	"main/domain"
	"main/model"
	"main/util"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const maxConnection = 50

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))

	logger.Default = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	second := 5
	for {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: false})
		if err != nil {
			log.Printf("Failed connect database. reconnect in %d seconds", second)
			time.Sleep(time.Duration(second) * time.Second)
			continue
		}
		sql, err := db.DB()
		sql.SetMaxOpenConns(maxConnection)
		sql.SetMaxIdleConns(maxConnection)
		sql.SetConnMaxIdleTime(time.Second)
		sql.SetConnMaxLifetime(time.Second)

		return db
	}
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) domain.IRepository {
	return &Repo{db: db}
}

func (r *Repo) GetCrawlableDate() time.Time {
	raw := model.Raw{}
	err := r.db.Last(&raw).Error
	if err != nil {
		return time.Date(2004, 2, 11, 0, 0, 0, 0, time.Local)
	}

	return raw.Date.Add(24 * time.Hour)
}

func (r *Repo) Insert(obj interface{}) error {
	return r.db.Create(obj).Error
}

func (r *Repo) Create(table string, obj interface{}) error {
	return r.db.Table(table).Create(obj).Error
}

func (r *Repo) AutoMigrate(objects ...interface{}) error {
	return r.db.AutoMigrate(objects...)
}

func (r *Repo) Migrate(table string, obj interface{}) error {
	return r.db.Table(table).AutoMigrate(obj)
}

func (r *Repo) GetConvertableDate() (time.Time, bool) {
	open := model.OpenDays{}
	raw := model.Raw{}

	err := r.db.Last(&open).Error
	if err != nil {

		err := r.db.First(&raw).Error
		if err != nil {
			return time.Time{}, false
		}
		return raw.Date, true

	}

	next := open.Date.Add(24 * time.Hour)

	err = r.db.First(&raw, "date = ?", next).Error
	if err != nil {
		return time.Time{}, false
	}
	return next, true
}

/* Goroutine handled */
func (r *Repo) GetRaw(date time.Time) ([]byte, error) {
	var raw model.Raw
	err := r.db.First(&raw, "date = ?", date).Error
	return raw.Body, err
}

func (r *Repo) GetStockHash() map[string]string {
	var stockList []*model.StockList
	hash := map[string]string{}

	err := r.db.Find(&stockList).Error
	if err != nil {
		log.Println("Error arr return")
		return hash
	}
	for _, s := range stockList {
		hash[s.StockID] = s.StockName
	}
	return hash
}

func (r *Repo) GetLastOpenDay() (time.Time, error) {
	open := model.OpenDays{}
	err := r.db.Where("state = ?", true).Last(&open).Error
	return open.Date, err
}

func (r *Repo) GetStock(id string) (model.Stock, error) {
	stock := model.NewStock()
	err := r.db.Model(&model.StockList{}).First(&stock, "stock_id = ?", id).Error
	if err != nil {
		return stock, fmt.Errorf("Failed to find stock")
	}

	var deals []model.Deal
	err = r.db.Table(util.StockTable(id)).Find(&deals).Error
	if err != nil {
		return stock, fmt.Errorf("Empty deals")
	}

	for _, d := range deals {
		stock.Deals[d.Date] = d
	}

	return stock, nil
}
func (r *Repo) GetStocksToday() ([]model.Stock, error) {
	var err error
	rawJson := model.RawJson{}

	date, err := r.GetLastOpenDay()
	if err != nil {
		return []model.Stock{}, err
	}

	raw, err := r.GetRaw(date)

	err = json.Unmarshal(raw, &rawJson)
	if err != nil {
		return []model.Stock{}, err
	}

	stocks := make([]model.Stock, 0, len(rawJson.Data9))
	for _, d := range rawJson.Data9 {
		deal := model.Deal{
			Date:        date,
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
		stock := model.NewStock()
		stock.ID = d[0]
		stock.Name = d[1]
		stock.Deals[deal.Date] = deal

		stocks = append(stocks, stock)
	}
	return stocks, nil
}
