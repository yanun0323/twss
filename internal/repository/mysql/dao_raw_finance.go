package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

var (
	_FinanceBeginPrevDate = time.Date(2005, time.September, 1, 0, 0, 0, 0, time.UTC)
)

func (dao MysqlDao) ListRawFinance(ctx context.Context, from, to time.Time) ([]model.RawFinance, error) {
	raws := []model.RawFinance{}
	err := dao.GetDriver(ctx).Where("date >= ? AND date <= ?", from, to).Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}

func (dao MysqlDao) GetRawFinanceDate(ctx context.Context, begin bool) (time.Time, error) {
	if begin {
		return _FinanceBeginPrevDate, nil
	}
	raw := model.RawFinance{}
	err := dao.GetDriver(ctx).Select("date").Last(&raw).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return _FinanceBeginPrevDate, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return raw.Date, nil
}

func (dao MysqlDao) GetRawFinance(ctx context.Context, date time.Time) (model.RawFinance, error) {
	raw := model.RawFinance{}
	err := dao.GetDriver(ctx).Table(raw.TableName()).Where("date = ?", date).Take(&raw).Error
	if err != nil {
		return model.RawFinance{}, err
	}
	return raw, nil
}

func (dao MysqlDao) InsertRawFinance(ctx context.Context, raw model.RawFinance) error {
	err := dao.GetDriver(ctx).Create(raw).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
