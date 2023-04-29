package mysql

import (
	"context"
	"stocker/internal/model"
	"time"
)

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
