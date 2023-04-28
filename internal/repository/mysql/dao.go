package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"stocker/internal/model"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
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
	_txKey                    = struct{}{}
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
	if os.Getenv("MODE") != "debug" {
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
		model.TradeDate{},
		model.RawEps{},
		model.RawTrade{},
		model.Stock{},
	)
}

func (dao MysqlDao) Migrate(table string, dst interface{}) {
	_ = dao.db.Table(table).AutoMigrate(dst)
}

func (dao MysqlDao) Debug(ctx context.Context) *gorm.DB {
	return dao.GetDriver(ctx)
}

func (dao MysqlDao) Tx(ctx context.Context, fc func(txCtx context.Context) error) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		_, ok := ctx.Value(_txKey).(*gorm.DB)
		if ok {
			return errors.New("transaction already exist")
		}

		txCtx := context.WithValue(ctx, _txKey, tx)
		return fc(txCtx)
	})
}

func (dao MysqlDao) IsErrRecordNotFound(err error) bool {
	return errors.Is(gorm.ErrRecordNotFound, err)
}

func (dao MysqlDao) GetDefaultStartDate() (time.Time, error) {
	return _defaultStartPreviousDate.Add(24 * time.Hour), nil
}

func (dao MysqlDao) GetDriver(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(_txKey).(*gorm.DB)
	if ok {
		return db
	}

	return dao.db
}

func isNotDuplicateEntryErr(err error) bool {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return true
	}
	return sqlErr.Number != 1062
}
