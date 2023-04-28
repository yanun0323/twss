package mysql

import (
	"context"
	"stocker/internal/model"
	"time"

	"github.com/yanun0323/pkg/logs"
)

func (dao MysqlDao) CheckTradeDate(ctx context.Context, date time.Time) error {
	return dao.GetDriver(ctx).Table(model.TradeDate{}.TableName()).Where("date = ?", date).Error
}

func (dao MysqlDao) GetLastTradeDate(ctx context.Context) (time.Time, error) {
	tradeDate := model.TradeDate{}
	if dao.GetDriver(ctx).Select("date").Last(&tradeDate).Error == nil {
		logs.Get(dao.ctx).Debug(tradeDate.Date)
		return tradeDate.Date, nil
	}
	return _defaultStartPreviousDate, nil
}

func (dao MysqlDao) InsertTradeDate(ctx context.Context, tradeDate model.TradeDate) error {
	err := dao.GetDriver(ctx).Create(tradeDate).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
