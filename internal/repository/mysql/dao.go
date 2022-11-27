package mysql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"stocker/internal/model"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_MAX_CONNECTION          = 50
	_RECONNECT_TIME_INTERVAL = 5 * time.Second
)

var (
	_defaultStartPreviousDate = time.Date(2004, time.February, 10, 0, 0, 0, 0, time.Local)
)

type MysqlDao struct {
	db  *gorm.DB
	ctx context.Context
}

func New(ctx context.Context) MysqlDao {
	dao := MysqlDao{
		db:  connectDB(ctx),
		ctx: ctx,
	}
	dao.AutoMigrate()
	return dao
}

func connectDB(ctx context.Context) *gorm.DB {
	l := logs.Get(ctx)
	loggers := logger.Default
	if os.Getenv("MODE") == "" {
		loggers = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))

	for {
		db, err := gorm.Open(sql.Open(dsn), &gorm.Config{
			Logger:                 loggers,
			SkipDefaultTransaction: false,
		})
		if err != nil {
			l.Warnf("connect database failed. reconnect in %d seconds, %+v", _RECONNECT_TIME_INTERVAL, err)
			time.Sleep(_RECONNECT_TIME_INTERVAL)
			continue
		}
		sql, err := db.DB()
		if err != nil {
			l.Warnf("connect database failed. reconnect in %d seconds, %+v", _RECONNECT_TIME_INTERVAL, err)
			time.Sleep(_RECONNECT_TIME_INTERVAL)
			continue
		}
		sql.SetMaxOpenConns(_MAX_CONNECTION)
		sql.SetMaxIdleConns(_MAX_CONNECTION)
		sql.SetConnMaxIdleTime(time.Second)
		sql.SetConnMaxLifetime(time.Second)

		return db
	}
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

func (dao MysqlDao) Debug() *gorm.DB {
	return dao.db
}

func (dao MysqlDao) ErrRecordNotFound() error {
	return gorm.ErrRecordNotFound
}

func (dao MysqlDao) CheckOpen(date time.Time) error {
	return dao.db.Table(model.Open{}.TableName()).Where("date = ?", date).Error
}

func (dao MysqlDao) CheckStock(date time.Time) error {
	table := model.DailyStock{ID: "2330"}.GetTableName()
	return dao.db.Table(table).Where("date = ?", date).Error
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
	open := model.Open{}
	if dao.db.Select("date").Last(&open).Error == nil {
		logs.Get(dao.ctx).Debug(open.Date)
		return open.Date, nil
	}
	return _defaultStartPreviousDate, nil
}

func (dao MysqlDao) GetStockMap() (model.StockMap, error) {
	list := model.StockList{}
	err := dao.db.Table(list.TableName()).Find(&list).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return model.StockMap{}, nil
	}
	if err != nil {
		return nil, err
	}
	return list.Map(), nil
}

func (dao MysqlDao) GetStock(id string) (model.Stock, error) {
	info := model.StockInfo{}
	data := []model.DailyStock{}

	if err := dao.db.Where("id = ?", id).Take(&info).Error; err != nil {
		return model.Stock{}, err
	}

	if err := dao.db.Table(model.DailyStock{ID: id}.GetTableName()).Find(&data).Error; err != nil {
		return model.Stock{}, err
	}

	return model.Stock{
		ID:        info.ID,
		Name:      info.Name,
		FirstDate: info.FirstDate,
		LastDate:  info.LastDate,
		Unable:    info.Unable,
		Trading:   data,
	}, nil
}

func (dao MysqlDao) GetDefaultStartDate() (time.Time, error) {
	return _defaultStartPreviousDate.Add(24 * time.Hour), nil
}

func (dao MysqlDao) GetLastDailyRawDate() (time.Time, error) {
	raw := model.DailyRaw{}
	err := dao.db.Select("date").Last(&raw).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return _defaultStartPreviousDate, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return raw.Date, nil
}

func (dao MysqlDao) GetDailyRaw(date time.Time) (model.DailyRaw, error) {
	raw := model.DailyRaw{}
	err := dao.db.Where("date = ?", date).Take(&raw).Error
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

func (dao MysqlDao) InsertDailyStockData(stock model.DailyStock) error {
	table := stock.GetTableName()
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
