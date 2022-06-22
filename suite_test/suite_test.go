package suite

import (
	"fmt"
	"main/config"
	"main/domain"
	"main/model"
	"main/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SVCSuite struct {
	suite.Suite
	repo domain.IRepository
}

func (su *SVCSuite) SetupSuite() {
	fmt.Println("SetupSuite")
	config.Init(".", "config")
	db := repository.ConnectDB()
	su.repo = repository.NewRepo(db)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SVCSuite))
}

func (su *SVCSuite) TestCreate() {
	su.repo.AutoMigrate(&model.TestStruct{})
	obj := model.TestStruct{Date: time.Date(2022, 3, 21, 0, 0, 0, 0, time.Local)}
	obj.State = true
	err := su.repo.Insert(obj)
	su.Nil(err, err)
}
