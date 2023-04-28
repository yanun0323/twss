package model

import "time"

type RawEps struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Body []byte    `gorm:"column:body"`
}

func (RawEps) TableName() string {
	return "raw_eps"
}
