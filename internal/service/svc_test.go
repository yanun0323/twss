package service

import (
	"context"
	"stocker/internal/repository"
	"stocker/pkg/infra"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	svc Service
	ctx context.Context
}

func (su *ServiceTestSuite) SetupTest() {
	su.Require().Nil(infra.Init("config-test"))
	su.ctx = context.Background()
	repo, err := repository.New(su.ctx)
	su.Require().NoError(err)
	su.svc, err = New(su.ctx, repo)
	su.Require().NoError(err)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (su *ServiceTestSuite) TestCrawl() {
	date, err := time.ParseInLocation("20060102", "20200501", time.Local)
	su.Require().Nil(err)
	su.Assert().Nil(su.svc.crawlRaw(date, CrawlEps))
}
