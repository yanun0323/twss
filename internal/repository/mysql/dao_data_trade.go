package mysql

import (
	"context"
	"stocker/internal/model"
	"time"
)

// Trade

func (dao MysqlDao) IsTradeExist(ctx context.Context, date time.Time) (bool, error) {
	table := model.Trade{ID: "2330"}.GetTableName()
	var count int64
	err := dao.GetDriver(ctx).Table(table).Where("date = ?", date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (dao MysqlDao) ListTrade(ctx context.Context, id string, from, to time.Time) ([]model.Trade, error) {
	table := model.Trade{ID: id}.GetTableName()
	var trades []model.Trade
	err := dao.GetDriver(ctx).Table(table).Where("date >= ? AND date <= ?", from, to).Find(&trades).Error
	if err != nil {
		return nil, err
	}
	return trades, nil
}

func (dao MysqlDao) GetTrade(ctx context.Context, id string, date time.Time) (model.Trade, error) {
	trade := model.Trade{ID: id}
	err := dao.GetDriver(ctx).Table(trade.GetTableName()).Where("date = ?", date).First(&trade).Error
	if err != nil {
		return model.Trade{}, err
	}
	return trade, nil
}

func (dao MysqlDao) InsertTrade(ctx context.Context, trade model.Trade) error {
	table := trade.GetTableName()
	dao.Migrate(table, trade)

	err := dao.GetDriver(ctx).Table(table).Create(trade).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

// Trade Date

func (dao MysqlDao) IsTradeDateExist(ctx context.Context, date time.Time) (bool, error) {
	var count int64
	err := dao.GetDriver(ctx).Table(model.TradeDate{}.TableName()).Where("date = ?", date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (dao MysqlDao) GetTradeDate(ctx context.Context, date time.Time) (model.TradeDate, error) {
	tradeDate := model.TradeDate{}
	err := dao.GetDriver(ctx).Where("date = ?", date).First(&tradeDate).Error
	if err != nil {
		return model.TradeDate{}, err
	}
	return tradeDate, nil
}

func (dao MysqlDao) GetLastTradeDate(ctx context.Context) (time.Time, error) {
	tradeDate := model.TradeDate{}
	err := dao.GetDriver(ctx).Select("date").Last(&tradeDate).Error
	if isNotFound(err) {
		return _TradeBeginPrevDate, nil
	}

	if err != nil {
		return _TradeBeginPrevDate, err
	}

	return tradeDate.Date, nil
}

func (dao MysqlDao) InsertTradeDate(ctx context.Context, tradeDate model.TradeDate) error {
	err := dao.GetDriver(ctx).Create(tradeDate).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
