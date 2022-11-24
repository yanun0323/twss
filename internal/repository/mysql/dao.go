package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"
	"time"

	"github.com/yanun0323/pkg/logs"
	"gorm.io/gorm"
)

type MysqlDao struct {
	db  *gorm.DB
	ctx context.Context
}

func New(ctx context.Context, db *gorm.DB) MysqlDao {
	dao := MysqlDao{
		db:  db,
		ctx: ctx,
	}
	dao.AutoMigrate()
	return dao
}

func (dao MysqlDao) AutoMigrate() {
	_ = dao.db.AutoMigrate(model.DailyRaw{})
}

func (dao MysqlDao) Migrate(table string, dst interface{}) {
	_ = dao.db.Table(table).AutoMigrate(dst)
}

func (dao MysqlDao) ListDailyRaws(from, to time.Time) ([]model.DailyRaw, error) {
	raws := []model.DailyRaw{}
	result := dao.db.Where("date >= ? AND date <= ?", from, to).Find(&raws)
	if result.Error != nil {
		return nil, result.Error
	}
	return raws, nil
}

func (dao MysqlDao) ListAllDailyRaws() ([]model.DailyRaw, error) {
	raws := []model.DailyRaw{}
	result := dao.db.Find(&raws)
	if result.Error != nil {
		return nil, result.Error
	}
	return raws, nil
}

func (dao MysqlDao) GetDailyRaw(date time.Time) (model.DailyRaw, error) {
	raw := model.DailyRaw{}
	result := dao.db.First(&raw, date)
	if result.Error != nil {
		return model.DailyRaw{}, result.Error
	}
	return raw, nil
}

func (dao MysqlDao) GetLastDailyRaw() (model.DailyRaw, error) {
	raw := model.DailyRaw{}
	result := dao.db.Last(&raw)
	if result.Error != nil {
		return model.DailyRaw{}, result.Error
	}
	return raw, nil
}

func (dao MysqlDao) InsertDailyRaw(raw model.DailyRaw) error {
	err := dao.db.Where("date = ?", raw.Date).Error
	if !errors.Is(gorm.ErrRecordNotFound, err) {
		logs.Get(dao.ctx).Debug("insert daily raw, data exist")
		return nil
	}

	if result := dao.db.Create(raw); result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao MysqlDao) InsertDailyStock(stock model.DailyStock) error {
	table := stock.TableName()
	dao.Migrate(table, stock)

	err := dao.db.Where("date = ?", stock.Date).Error
	if !errors.Is(gorm.ErrRecordNotFound, err) {
		logs.Get(dao.ctx).Debug("insert daily stock, data exist")
		return nil
	}

	result := dao.db.Table(table).Create(stock)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
