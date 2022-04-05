package model

import "time"

/*
	[0]  "證券代號"
	[1]  "證券名稱"
	[2]  "成交股數"
	[3]  "成交筆數"
	[4]  "成交金額"
	[5]  "開盤價"
	[6]  "最高價"
	[7]  "最低價"
	[8]  "收盤價"
	[9]  "漲跌(+/-)"
	[10] "漲跌價差"
	[11] "最後揭示買價"
	[12] "最後揭示買量"
	[13] "最後揭示賣價"
	[14] "最後揭示賣量"
	[15] "本益比"
*/
type Stock struct {
	ID    string `gorm:"column:stock_id;primaryKey;not null;size:30"`
	Name  string `gorm:"column:stock_name;not null;size:30"`
	Deals map[time.Time]Deal
}

func NewStock() Stock {
	return Stock{
		Deals: map[time.Time]Deal{},
	}
}
