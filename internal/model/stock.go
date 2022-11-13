package model

import "time"

type Raw struct {
	Date time.Time `gorm:"column:date;primaryKey;not null"`
	Body string    `gorm:"column:body"`
}

func (Raw) TableName() string {
	return "price_raw_body"
}

type RawDaily struct {
	Stat   string     `json:"stat"`
	Date   string     `json:"date"`
	Title  string     `json:"title"`
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data"`
	Notes  []string   `json:"notes"`
}

func (raw RawDaily) ParseStocks() []StockDaily {
	sd := make([]StockDaily, 0, len(raw.Data))
	for _, s := range raw.Data {
		sd = append(sd, StockDaily{
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

// "證券代號","證券名稱","成交股數","成交金額","開盤價","最高價","最低價","收盤價","漲跌價差","成交筆數"
type StockDaily struct {
	ID, Name                                         string
	TradeShare, TradeMoney                           string
	PriceOpen, PriceLowest, PriceHighest, PriceClose string
	TradeGrade, TradeCount                           string
}
