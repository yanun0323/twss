package model

import "time"

type Raw struct {
	Date time.Time `gorm:"column:date;primaryKey;not null"`
	Body []byte    `gorm:"column:body"`
}

func (r Raw) TableName() string {
	return "price_raw_body"
}
