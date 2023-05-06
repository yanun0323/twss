package service

import (
	"context"
	"stocker/internal/model"
	"stocker/internal/repository"
	"stocker/pkg/infra"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	svc Service
	ctx context.Context
}

func (su *ServiceTestSuite) SetupTest() {
	su.Require().Nil(infra.Init("config"))
	su.ctx = context.Background()
	repo, err := repository.New(su.ctx)
	su.Require().NoError(err)
	su.svc, err = New(su.ctx, repo)
	su.Require().NoError(err)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (su *ServiceTestSuite) TestPower() {
	in := model.PowerInput{
		ID:   "3380",
		From: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2023, 4, 28, 0, 0, 0, 0, time.UTC),
	}
	result, err := su.svc.CalculatePower(in)
	su.Require().NoError(err, err)

	// su.svc.Log.Infof("ID: %s", result.ID)
	base := decimal.Zero
	for d := in.From; d.Before(in.To); d = d.Add(24 * time.Hour) {
		power, ok := result.Power[d]
		if !ok {
			continue
		}
		if base.IsZero() {
			base = power
		}

		danger := power.Sub(base).DivRound(base, 5)
		dS := danger.String()
		if danger.Sign() != -1 {
			dS = "+" + dS
		}
		// su.svc.Log.Infof("[%s] D: %s, P: %s", d.Format("2006-01-02"), dS, power)
	}
}
