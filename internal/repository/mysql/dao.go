package mysql

import (
	"stocker/internal/model"

	"gorm.io/gorm"
)

type MysqlDao struct {
	db *gorm.DB
}

func New(db *gorm.DB) MysqlDao {
	dao := MysqlDao{
		db: db,
	}
	dao.Migrate()
	return dao
}

func (dao MysqlDao) Migrate() {
	_ = dao.db.AutoMigrate(model.Raw{})
}

func (dao MysqlDao) GetLastRaw() (model.Raw, error) {
	raw := model.Raw{}
	result := dao.db.Table(raw.TableName()).Last(&raw)
	if result.Error != nil {
		return model.Raw{}, result.Error
	}
	return raw, nil
}

func (dao MysqlDao) InsertRaw(raw model.Raw) error {
	result := dao.db.Table(raw.TableName()).Create(raw)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
