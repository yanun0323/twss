package model

import (
	"time"
)

type Deal struct {
	Date        time.Time `gorm:"column:date;primaryKey"`
	Volume      string    `gorm:"size:255"`
	VolumeMoney string    `gorm:"size:255"`
	Start       string    `gorm:"size:255"`
	Max         string    `gorm:"size:255"`
	Min         string    `gorm:"size:255"`
	End         string    `gorm:"size:255"`
	Spread      string    `gorm:"size:30"`
	Per         string    `gorm:"size:255"`
}
