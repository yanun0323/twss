package model

import (
	"context"
	"stocker/internal/util"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/yanun0323/pkg/logs"
)

// RawFinance 爬蟲的每日財務指標
type RawFinance struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Body []byte    `gorm:"column:body"`
}

func (RawFinance) TableName() string {
	return "raw_finance"
}

func (raw RawFinance) GetData() (RawData, error) {
	data := RawFinanceData{}
	if err := json.Unmarshal([]byte(raw.Body), &data); err != nil {
		return RawFinanceData{}, err
	}
	data.Date = raw.Date
	return data, nil
}

// RawFinanceData 爬蟲的每日財務指標解析後資料
type RawFinanceData struct {
	Date   time.Time       `gorm:"column:date;primaryKey" json:"-"`
	Stat   string          `json:"stat,omitempty"`
	Title  string          `json:"title,omitempty"`
	Fields []string        `json:"fields,omitempty"`
	Data   [][]interface{} `json:"data,omitempty"`
}

func (raw RawFinanceData) IsOK() bool {
	return raw.Stat == "OK" && len(raw.Data) != 0
}

// Parse 解析資料，分為舊資料(5欄)與新資料(7欄)
func (raw RawFinanceData) Parse() []interface{} {
	sd := make([]interface{}, 0, len(raw.Data))

	for _, s := range raw.Data {
		switch len(s) {
		case 5:
			sd = append(sd, raw.ParseOld(raw.Date, s))
		case 7:
			sd = append(sd, raw.ParseNew(raw.Date, s))
		default:
			logs.Get(context.Background()).Errorf("unknown finance data: %s - %s", util.LogDate(raw.Date), s[0].(string))
		}
	}
	return sd
}

// [0:證券代號 1:證券名稱 2:本益比 3:殖利率(%) 4:股價淨值比]
func (raw RawFinanceData) ParseOld(date time.Time, s []interface{}) Finance {
	return Finance{
		Date:   date,
		ID:     s[0].(string),
		Name:   s[1].(string),
		DY:     util.Decimal(s[3].(string)).Mul(_hundred),
		PER:    util.Decimal(s[2].(string)),
		PBR:    util.Decimal(s[4].(string)),
		Year:   date.Year() - 1911,
		Season: -1,
	}
}

// [0:證券代號 1:證券名稱 2:殖利率(%) 3:股利年度(Int) 4:本益比 5:股價淨值比 6:財報年/季]
func (raw RawFinanceData) ParseNew(date time.Time, s []interface{}) Finance {
	year, session := parseYearSeason(s[6].(string))
	return Finance{
		Date:   date,
		ID:     s[0].(string),
		Name:   s[1].(string),
		DY:     util.Decimal(s[2].(string)).Mul(_hundred),
		PER:    util.Decimal(s[4].(string)),
		PBR:    util.Decimal(s[5].(string)),
		Year:   year,
		Season: session,
	}
}

func parseYearSeason(s string) (int, int) {
	var year, session int
	ys := strings.Split(s, "/")
	if len(ys) == 2 {
		year = util.Int(ys[0])
		session = util.Int(ys[1])
	}
	return year, session
}
