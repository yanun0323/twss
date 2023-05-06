package mysql

import (
	"context"
	"stocker/internal/model"
	"time"
)

// Finance

func (dao MysqlDao) IsFinanceExist(ctx context.Context, date time.Time) (bool, error) {
	table := model.Finance{ID: "2330"}.GetTableName()
	var count int64
	err := dao.GetDriver(ctx).Table(table).Where("date = ?", date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (dao MysqlDao) InsertFinance(ctx context.Context, f model.Finance) error {
	table := f.GetTableName()
	dao.Migrate(table, f)

	err := dao.GetDriver(ctx).Table(table).Create(f).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

// Finance Date

func (dao MysqlDao) IsFinanceDateExist(ctx context.Context, date time.Time) (bool, error) {
	var count int64
	err := dao.GetDriver(ctx).Table(model.FinanceDate{}.TableName()).Where("date = ?", date).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (dao MysqlDao) GetFinanceDate(ctx context.Context, date time.Time) (model.FinanceDate, error) {
	fd := model.FinanceDate{}
	err := dao.GetDriver(ctx).Where("date = ?", date).First(&fd).Error
	if err != nil {
		return model.FinanceDate{}, err
	}
	return fd, nil
}

func (dao MysqlDao) GetLastFinanceDate(ctx context.Context) (time.Time, error) {
	fd := model.FinanceDate{}
	err := dao.GetDriver(ctx).Select("date").Last(&fd).Error
	if isNotFound(err) {
		return _FinanceBeginPrevDate, nil
	}

	if err != nil {
		return _FinanceBeginPrevDate, err
	}

	return fd.Date, nil
}

func (dao MysqlDao) InsertFinanceDate(ctx context.Context, fd model.FinanceDate) error {
	err := dao.GetDriver(ctx).Create(fd).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}
