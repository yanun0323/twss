package suite

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
