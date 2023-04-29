package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

var (
	_EpsBeginPrevDate = time.Date(2005, time.September, 1, 0, 0, 0, 0, time.UTC)
)

func (dao MysqlDao) ListRawEps(ctx context.Context, from, to time.Time) ([]model.RawEps, error) {
	raws := []model.RawEps{}
	err := dao.GetDriver(ctx).Where("date >= ? AND date <= ?", from, to).Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}

func (dao MysqlDao) GetRawEpsDate(ctx context.Context, begin bool) (time.Time, error) {
	if begin {
		return _EpsBeginPrevDate, nil
	}
	raw := model.RawEps{}
	err := dao.GetDriver(ctx).Select("date").Last(&raw).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return _EpsBeginPrevDate, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return raw.Date, nil
}

func (dao MysqlDao) GetRawEps(ctx context.Context, date time.Time) (model.RawEps, error) {
	raw := model.RawEps{}
	err := dao.GetDriver(ctx).Table(raw.TableName()).Where("date = ?", date).Take(&raw).Error
	if err != nil {
		return model.RawEps{}, err
	}
	return raw, nil
}

func (dao MysqlDao) InsertRawEps(ctx context.Context, raw model.RawEps) error {
	err := dao.GetDriver(ctx).Create(raw).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
