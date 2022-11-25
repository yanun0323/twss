package mysql

import (
	"context"
	"errors"
	"stocker/internal/model"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	_defaultStartDate = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.Local)
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
	_ = dao.db.AutoMigrate(
		model.Open{},
		model.DailyRaw{},
		model.StockInfo{},
	)
}

func (dao MysqlDao) Migrate(table string, dst interface{}) {
	_ = dao.db.Table(table).AutoMigrate(dst)
}

func (dao MysqlDao) ListAllDailyRaws() ([]model.DailyRaw, error) {
	raws := []model.DailyRaw{}
	err := dao.db.Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}

func (dao MysqlDao) ListDailyRaws(from, to time.Time) ([]model.DailyRaw, error) {
	raws := []model.DailyRaw{}
	err := dao.db.Where("date >= ? AND date <= ?", from, to).Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}

func (dao MysqlDao) GetLastOpenDate() (time.Time, error) {
	t := time.Time{}
	err := dao.db.Select("date").Last(&t).Error
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (dao MysqlDao) GetStockMap() (model.StockMap, error) {
	list := model.StockList{}
	err := dao.db.Table(list.TableName()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list.Map(), nil
}

func (dao MysqlDao) GetLastDailyRawDate() (time.Time, error) {
	t := time.Time{}
	err := dao.db.Select("date").Last(&t).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return _defaultStartDate, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (dao MysqlDao) GetDailyRaw(date time.Time) (model.DailyRaw, error) {
	raw := model.DailyRaw{}
	err := dao.db.First(&raw, date).Error
	if err != nil {
		return model.DailyRaw{}, err
	}
	return raw, nil
}

func (dao MysqlDao) InsertOpen(open model.Open) error {
	err := dao.db.Create(open).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

func (dao MysqlDao) InsertDailyRaw(raw model.DailyRaw) error {
	err := dao.db.Create(raw).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

func (dao MysqlDao) InsertStockList(info model.StockInfo) error {
	err := dao.db.Create(info).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

func (dao MysqlDao) InsertDailyStock(stock model.DailyStock) error {
	table := stock.TableName()
	dao.Migrate(table, stock)

	err := dao.db.Table(table).Create(stock).Error
	if err != nil && isNotDuplicateEntryErr(err) {
		return err
	}
	return nil
}

func isNotDuplicateEntryErr(err error) bool {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return true
	}
	return sqlErr.Number != 1062
}
