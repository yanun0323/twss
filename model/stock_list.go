package model

type StockList struct {
	StockID   string `gorm:"column:stock_id;primaryKey;not null;size:30"`
	StockName string `gorm:"column:stock_name;not null;size:30"`
}

func (s StockList) TableName() string {
	return "stock_list"
}
