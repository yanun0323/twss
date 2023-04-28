package model

import (
	"regexp"
	"stocker/internal/util"
	"time"

	"github.com/goccy/go-json"

	"github.com/shopspring/decimal"
)

var (
	_regexpTradeSymbol = regexp.MustCompile(".</p>$")
)

// RawTrade 爬蟲的每日交易
type RawTrade struct {
	Date time.Time `gorm:"column:date;primaryKey"`
	Body []byte    `gorm:"column:body"`
}

func (RawTrade) TableName() string {
	return "raw_trade"
}

func (raw RawTrade) GetData() (RawTradeData, error) {
	data := RawTradeData{}
	if err := json.Unmarshal([]byte(raw.Body), &data); err != nil {
		return RawTradeData{}, err
	}
	data.Date = raw.Date
	return data, nil
}

// RawTradeData 爬蟲的每日交易解析後資料
type RawTradeData struct {
	Date      time.Time       `gorm:"column:date;primaryKey" json:"-"`
	Stat      string          `json:"stat,omitempty"`
	Title     string          `json:"title,omitempty"`
	Fields1   []string        `json:"fields1,omitempty"`
	Fields2   []string        `json:"fields2,omitempty"`
	Fields3   []string        `json:"fields3,omitempty"`
	Fields4   []string        `json:"fields4,omitempty"`
	Fields5   []string        `json:"fields5,omitempty"`
	Fields6   []string        `json:"fields6,omitempty"`
	Fields7   []string        `json:"fields7,omitempty"`
	Fields8   []string        `json:"fields8,omitempty"`
	Fields9   []string        `json:"fields9,omitempty"`
	Subtitle1 string          `json:"subtitle1,omitempty"`
	Subtitle2 string          `json:"subtitle2,omitempty"`
	Subtitle3 string          `json:"subtitle3,omitempty"`
	Subtitle4 string          `json:"subtitle4,omitempty"`
	Subtitle5 string          `json:"subtitle5,omitempty"`
	Subtitle6 string          `json:"subtitle6,omitempty"`
	Subtitle7 string          `json:"subtitle7,omitempty"`
	Subtitle8 string          `json:"subtitle8,omitempty"`
	Subtitle9 string          `json:"subtitle9,omitempty"`
	Data1     [][]string      `json:"data1,omitempty"`
	Data2     [][]string      `json:"data2,omitempty"`
	Data3     [][]string      `json:"data3,omitempty"`
	Data4     [][]string      `json:"data4,omitempty"`
	Data5     [][]string      `json:"data5,omitempty"`
	Data6     [][]string      `json:"data6,omitempty"`
	Data7     [][]interface{} `json:"data7,omitempty"`
	Data8     [][]string      `json:"data8,omitempty"`
	Data9     [][]string      `json:"data9,omitempty"`
	Notes     []string        `json:"notes,omitempty"`
}

func (raw *RawTradeData) StockData() [][]string {
	// data8: before 2011/7/31
	// data9: '2006/09/29' and after 2011/7/31
	if len(raw.Data9) > len(raw.Data8) {
		return raw.Data9
	}
	return raw.Data8
}

// [0:證券代號 1:證券名稱 2:成交股數 3:成交筆數 4:成交金額 5:開盤價 6:最高價 7:最低價 8:收盤價 9:漲跌(+/-) 10:漲跌價差 11:最後揭示買價 12:最後揭示買量 13:最後揭示賣價 14:最後揭示賣量 15:本益比]
func (raw *RawTradeData) ParseStock() []Stock {
	sd := make([]Stock, 0, len(raw.StockData()))
	for _, s := range raw.StockData() {
		symbol := parseSymbol(s[9])
		grade := s[10]
		if symbol != " " {
			grade = symbol + grade
		}
		sd = append(sd, Stock{
			Date:         raw.Date,
			ID:           s[0],
			Name:         s[1],
			TradeShare:   util.Decimal(s[2]),
			TradeCount:   util.Decimal(s[3]),
			TradeMoney:   util.Decimal(s[4]),
			PriceOpen:    util.Decimal(s[5]),
			PriceHighest: util.Decimal(s[6]),
			PriceLowest:  util.Decimal(s[7]),
			PriceClose:   util.Decimal(s[8]),
			TradeSymbol:  symbol,
			TradeGrade:   util.Decimal(grade),
			Percentage:   calculatePercentage(grade, s[5]),
			Limit:        calculateLimit(s[8], grade),
		})
	}
	return sd
}

func parseSymbol(s string) string {
	return string(_regexpTradeSymbol.FindString(s)[0])
}

func calculatePercentage(gradeStr, openStr string) decimal.Decimal {
	open, err := decimal.NewFromString(openStr)
	if err != nil || open.IsZero() {
		return decimal.Zero
	}
	grade, err := decimal.NewFromString(gradeStr)
	if err != nil {
		return decimal.Zero
	}
	return grade.Div(open).Shift(2).Round(2).Truncate(2)
}

func calculateLimit(closeStr, gradeStr string) bool {
	close, err := decimal.NewFromString(closeStr)
	if err != nil {
		return false
	}
	grade, err := decimal.NewFromString(gradeStr)
	if err != nil || grade.IsZero() {
		return false
	}

	begin := close.Sub(grade)
	_, fixed := priceThreshold(begin)
	interval := begin.Shift(-1).RoundFloor(fixed)
	if grade.IsNegative() {
		interval = interval.Neg()
	}
	expected := fixPrice(begin.Add(interval))
	if begin.Sign() > 0 {
		return close.GreaterThanOrEqual(expected)
	}
	return close.LessThanOrEqual(expected)
}

func fixPrice(price decimal.Decimal) decimal.Decimal {
	if price.LessThanOrEqual(decimal.Zero) {
		return price
	}
	threshold, fixed := priceThreshold(price)
	price = price.Round(fixed)
	prefix := price.Shift(-1).Truncate(fixed).Shift(1)
	suffix := price.Sub(prefix).Div(threshold).Truncate(0).Mul(threshold)
	if suffix.IsNegative() {
		suffix = decimal.Zero
	}
	return prefix.Add(suffix)
}

var (
	_ten                  = decimal.New(10, 0)
	_tenThreshold         = decimal.New(1, -2)
	_fifty                = decimal.New(50, 0)
	_fiftyThreshold       = decimal.New(5, -2)
	_hundred              = decimal.New(100, 0)
	_hundredThreshold     = decimal.New(1, -1)
	_fiveHundred          = decimal.New(500, 0)
	_fiveHundredThreshold = decimal.New(5, -1)
	_thousand             = decimal.New(1000, 0)
	_thousandThreshold    = decimal.New(1, 0)
	_otherThreshold       = decimal.New(5, 0)
)

func priceThreshold(price decimal.Decimal) (decimal.Decimal, int32) {
	price = price.Round(2)
	if price.LessThan(_ten) {
		return _tenThreshold, 2
	}
	if price.LessThan(_fifty) {
		return _fiftyThreshold, 2
	}
	price = price.Round(1)
	if price.LessThan(_hundred) {
		return _hundredThreshold, 1
	}
	if price.LessThan(_fiveHundred) {
		return _fiveHundredThreshold, 1
	}
	price = price.Round(0)
	if price.LessThan(_thousand) {
		return _thousandThreshold, 0
	}
	return _otherThreshold, 0
}
