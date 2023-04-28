package model

import "time"

// RawTrade 爬蟲的每日EPS
type RawEps struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Body []byte    `gorm:"column:body"`
}

func (RawEps) TableName() string {
	return "raw_eps"
}
