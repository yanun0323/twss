package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"

	"gorm.io/gorm"
)

func (dao MysqlDao) ListStocks(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}
	err := dao.GetDriver(ctx).Table(model.Stock{}.TableName()).Find(&stocks).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (dao MysqlDao) InsertStock(ctx context.Context, s model.Stock) error {
	err := dao.GetDriver(ctx).Create(s).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
