package cosmos_test

import (
	"testing"

	"github.com/lyfeblocnetwork/os/ante/testutils"
	"github.com/stretchr/testify/suite"
)

type AnteTestSuite struct {
	*testutils.AnteTestSuite
}

func TestAnteTestSuite(t *testing.T) {
	baseSuite := new(testutils.AnteTestSuite)
	baseSuite.WithLondonHardForkEnabled(true)
	baseSuite.WithFeemarketEnabled(true)

	suite.Run(t, &AnteTestSuite{baseSuite})
}
