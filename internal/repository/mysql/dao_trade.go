package mysql

import (
	"context"
	"stocker/internal/model"
	"time"
)

func (dao MysqlDao) IsTradeExist(ctx context.Context, date time.Time) (bool, error) {
	table := model.Trade{ID: "2330"}.GetTableName()
	var count int64
	err := dao.GetDriver(ctx).Table(table).Where("date = ?", date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
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
