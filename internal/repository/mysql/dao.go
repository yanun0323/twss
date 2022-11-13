package mysql

import (
	"stocker/internal/model"
	"time"

	"gorm.io/gorm"
)

type MysqlDao struct {
	db *gorm.DB
}

func New(db *gorm.DB) MysqlDao {
	dao := MysqlDao{
		db: db,
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

func (dao MysqlDao) GetLastDailyRaw() (model.DailyRaw, error) {
	raw := model.DailyRaw{}
	result := dao.db.Last(&raw)
	if result.Error != nil {
		return model.DailyRaw{}, result.Error
	}
	return raw, nil
}

func (dao MysqlDao) InsertDailyRaw(raw model.DailyRaw) error {
	result := dao.db.Create(raw)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao MysqlDao) InsertDailyStock(stock model.DailyStock) error {
	table := stock.TableName()
	dao.Migrate(table, stock)

	result := dao.db.Table(table).Create(stock)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
