package repository

import (
	"fmt"
	"log"
	"main/domain"
	"os"
	"time"

	"main/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const maxConnection = 50

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))

	logger.Default = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	second := 5
	for {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: false})
		if err != nil {
			log.Printf("Failed connect database. reconnect in %d seconds", second)
			time.Sleep(time.Duration(second) * time.Second)
			continue
		}
		sql, err := db.DB()
		sql.SetMaxOpenConns(maxConnection)
		sql.SetMaxIdleConns(maxConnection)
		sql.SetConnMaxIdleTime(time.Second)
		sql.SetConnMaxLifetime(time.Second)

		return db
	}
}

type Repo struct {
	dao.MysqlDao
}

func NewRepo(db *gorm.DB) domain.IRepository {
	return &Repo{
		MysqlDao: dao.NewMysqlDao(db),
	}
}
