package model

import (
	"time"
)

type StockInfo struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	FirstDate time.Time `gorm:"column:first_date;not null"`
	LastDate  time.Time `gorm:"column:last_date;not null"`
	Unable    bool      `gorm:"column:unable;not null"`
}

func (StockInfo) TableName() string {
	return "stock_list"
}

type StockList []StockInfo

func (StockList) TableName() string {
	return "stock_list"
}

func (l StockList) Map() StockMap {
	m := StockMap{}
	for _, v := range l {
		m[v.ID] = v.Name
	}
	return m
}

type StockMap map[string]string

func (m StockMap) List() StockList {
	list := make(StockList, 0, len(m))
	for id, name := range m {
		list = append(list, StockInfo{
			ID:   id,
			Name: name,
		})
	}
	return list
}
