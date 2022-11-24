package repository

import (
	"context"
	"fmt"
	"stocker/internal/domain"
	"stocker/internal/repository/mysql"
	"time"

	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	_MAX_CONNECTION          = 30
	_RECONNECT_TIME_INTERVAL = 5 * time.Second
)

type Repo struct {
	mysql.MysqlDao
}

func New(ctx context.Context) domain.Repository {
	return &Repo{
		MysqlDao: mysql.New(ctx, connectDB(ctx)),
	}
}

func connectDB(ctx context.Context) *gorm.DB {
	l := logs.Get(ctx)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))

	for {
		db, err := gorm.Open(sql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: false})
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
