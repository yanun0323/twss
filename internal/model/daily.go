package model

import "time"

type DailyRaw struct {
	Date time.Time `gorm:"column:date;primaryKey;not null"`
	Body []byte    `gorm:"column:body"`
}

func (DailyRaw) TableName() string {
	return "daily_raw"
}

type DailyRawData struct {
	Stat   string     `json:"stat"`
	Date   string     `json:"date"`
	Title  string     `json:"title"`
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data9"`
	Notes  []string   `json:"notes"`
}

func (raw *DailyRawData) ParseStock(date time.Time) []DailyStock {
	sd := make([]DailyStock, 0, len(raw.Data))
	for _, s := range raw.Data {
		sd = append(sd, DailyStock{
			Date:         date,
			ID:           s[0],
			Name:         s[1],
			TradeShare:   s[2],
			TradeMoney:   s[3],
			PriceOpen:    s[4],
			PriceLowest:  s[5],
			PriceHighest: s[6],
			PriceClose:   s[7],
			TradeGrade:   s[8],
			TradeCount:   s[9],
		})
	}
	return sd
}
