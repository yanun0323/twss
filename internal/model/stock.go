package model

import "time"

type Stock struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	FirstDate time.Time    `json:"firstDate"`
	LastDate  time.Time    `json:"lastDate"`
	Unable    bool         `json:"unable"`
	Trading   []DailyStock `json:"trading"`
}
