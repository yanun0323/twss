package model

import (
	"fmt"
	"strings"
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

func NewDealFromDataString(d []string, date time.Time) Deal {
	prefix := ""
	if strings.Contains(d[9], "green") {
		prefix = "-"
	}
	if strings.Contains(d[9], "red") {
		prefix = "+"
	}
	return Deal{
		Date:        date,
		Volume:      d[2],
		VolumeMoney: d[3],
		Start:       d[5],
		Max:         d[6],
		Min:         d[7],
		End:         d[8],
		Spread:      fmt.Sprintf("%s%s", prefix, d[10]),
		Per:         d[15],
	}
}
