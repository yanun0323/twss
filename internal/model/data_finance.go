package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Finance 財務資料
type Finance struct {
	Date   time.Time       `gorm:"column:date;primaryKey" json:"date"`
	ID     string          `gorm:"-" json:"id,omitempty"`
	Name   string          `gorm:"-" json:"name,omitempty"`
	DY     decimal.Decimal `gorm:"column:dividend_yield;not null" json:"dividend_yield"` /* DY 殖利率 = (股息 / 股價) * 100% */
	PER    decimal.Decimal `gorm:"column:per;not null" json:"per"`                       /* PER 本益比 = 股價 / 每股盈餘 */
	PBR    decimal.Decimal `gorm:"column:pbr;not null" json:"pbr"`                       /* PBR 股價淨值比 = 股價 / 每股淨值 */
	Year   int             `gorm:"column:year;not null" json:"year"`                     /* 財報年 */
	Season int             `gorm:"column:season;not null" json:"season"`                 /* 財報季 */
}

func (f Finance) GetTableName() string {
	return "finance_" + f.ID
}

// FinanceDate 財務日期
type FinanceDate struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Open bool      `gorm:"column:is_open;not null"`
}

func (fd FinanceDate) IsOpen() bool {
	return fd.Open
}

func (FinanceDate) TableName() string {
	return "finance_date"
}
