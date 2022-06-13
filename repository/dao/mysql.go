package dao

import (
	"encoding/json"
	"fmt"
	"log"
	"main/model"
	"main/util"
	"time"

	"gorm.io/gorm"
)

type MysqlDao struct {
	db *gorm.DB
}

func NewMysqlDao(db *gorm.DB) MysqlDao {
	return MysqlDao{db: db}
}

func (r *MysqlDao) GetCrawlableDate(checkMode bool) time.Time {
	if checkMode {
		return time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local)
	}

	raw := model.Raw{}
	err := r.db.Last(&raw).Error
	if err != nil {
		return time.Date(2004, 2, 11, 0, 0, 0, 0, time.Local)
	}

	return raw.Date.Add(24 * time.Hour)
}

func (r *MysqlDao) Insert(obj interface{}) error {
	return r.db.Save(obj).Error
}

func (r *MysqlDao) InsertWithTableName(table string, obj interface{}) error {
	return r.db.Table(table).Save(obj).Error
}

// arg must be pointer
func (r *MysqlDao) AutoMigrate(obj ...interface{}) error {
	return r.db.AutoMigrate(obj...)
}

// arg must be pointer
func (r *MysqlDao) Migrate(table string, obj interface{}) error {
	return r.db.Table(table).AutoMigrate(obj)
}

func (r *MysqlDao) GetConvertibleDate(checkMode bool) (time.Time, error) {
	if checkMode {
		return time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local), nil
	}
	open := model.OpenDays{}

	err := r.db.Last(&open).Error
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get convertible date")
	}

	return open.Date.Add(24 * time.Hour), nil
}

/* Goroutine handled */
func (r *MysqlDao) GetRaw(date time.Time) ([]byte, error) {
	var raw model.Raw
	err := r.db.First(&raw, "date = ?", date).Error
	return raw.Body, err
}

func (r *MysqlDao) GetStockHash() map[string]string {
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

func (r *MysqlDao) GetLastOpenDay() (time.Time, error) {
	open := model.OpenDays{}
	err := r.db.Where("state = ?", true).Last(&open).Error
	return open.Date, err
}

func (r *MysqlDao) GetStock(id string) (model.Stock, error) {
	stock := model.NewStock()
	err := r.db.Model(&model.StockList{}).First(&stock, "stock_id = ?", id).Error
	if err != nil {
		return stock, fmt.Errorf("failed to find stock")
	}

	var deals []model.Deal
	err = r.db.Table(util.StockTable(id)).Find(&deals).Error
	if err != nil {
		return stock, fmt.Errorf("empty deals")
	}

	for _, d := range deals {
		stock.Deals[d.Date] = d
	}

	return stock, nil
}
func (r *MysqlDao) GetStocksToday() ([]model.Stock, error) {
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
