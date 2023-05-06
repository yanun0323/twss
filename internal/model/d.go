package model

/* Check Interface Implement */
var (
	_ Raw = (*RawTrade)(nil)
	_ Raw = (*RawFinance)(nil)

	_ RawData = (*RawTradeData)(nil)
	_ RawData = (*RawFinanceData)(nil)

	_ DataDate = (*TradeDate)(nil)
	_ DataDate = (*FinanceDate)(nil)
)

type Raw interface {
	GetData() (RawData, error)
}

type RawData interface {
	IsOK() bool
	Parse() []interface{}
}

type DataDate interface {
	IsOpen() bool
}
