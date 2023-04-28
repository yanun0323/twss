package model

import (
	"time"
)

const (
	_STOCK_LIST_TABLE_NAME = "stock_list"
)

type StockListUnit struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	FirstDate time.Time `gorm:"column:first_date;not null"`
	LastDate  time.Time `gorm:"column:last_date;not null"`
	Unable    bool      `gorm:"column:unable;not null"`
}

func (StockListUnit) TableName() string {
	return _STOCK_LIST_TABLE_NAME
}

type StockList []StockListUnit

func (StockList) TableName() string {
	return _STOCK_LIST_TABLE_NAME
}

func (l StockList) Map() StockMap {
	m := StockMap{}
	for _, v := range l {
		m[v.ID] = v.Name
	}
	return m
}
