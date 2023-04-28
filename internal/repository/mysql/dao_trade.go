package mysql

import (
	"context"
	"stocker/internal/model"
	"time"
)

func (dao MysqlDao) CheckTrade(ctx context.Context, date time.Time) error {
	table := model.Trade{ID: "2330"}.GetTableName()
	return dao.GetDriver(ctx).Table(table).Where("date = ?", date).Error
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
