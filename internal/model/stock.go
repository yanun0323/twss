package model

import "time"

type Stock struct {
	ID        string
	Name      string
	FirstDate time.Time
	LastDate  time.Time
	Unable    bool
	Trading   []DailyStockData
}
