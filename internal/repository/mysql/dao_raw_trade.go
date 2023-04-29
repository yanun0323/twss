package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

var (
	_TradeBeginPrevDate = time.Date(2004, time.February, 10, 0, 0, 0, 0, time.UTC)
)

func (dao MysqlDao) ListRawTrade(ctx context.Context, from, to time.Time) ([]model.RawTrade, error) {
	raws := []model.RawTrade{}
	err := dao.GetDriver(ctx).Where("date >= ? AND date <= ?", from, to).Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}

func (dao MysqlDao) GetRawTradeDate(ctx context.Context, begin bool) (time.Time, error) {
	if begin {
		return _TradeBeginPrevDate, nil
	}
	raw := model.RawTrade{}
	err := dao.GetDriver(ctx).Select("date").Last(&raw).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return _TradeBeginPrevDate, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return raw.Date, nil
}

func (dao MysqlDao) GetRawTrade(ctx context.Context, date time.Time) (model.RawTrade, error) {
	raw := model.RawTrade{}
	err := dao.GetDriver(ctx).Table(raw.TableName()).Where("date = ?", date).Take(&raw).Error
	if err != nil {
		return model.RawTrade{}, err
	}
	return raw, nil
}

func (dao MysqlDao) InsertRawTrade(ctx context.Context, raw model.RawTrade) error {
	err := dao.GetDriver(ctx).Create(raw).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
