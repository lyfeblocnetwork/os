package evm_test

import (
	"testing"

	"github.com/lyfeblocnetwork/os/ante/testutils"
	"github.com/stretchr/testify/suite"
)

type AnteTestSuite struct {
	*testutils.AnteTestSuite
	useLegacyEIP712TypedData bool
}

func TestAnteTestSuite(t *testing.T) {
	baseSuite := new(testutils.AnteTestSuite)
	baseSuite.WithLondonHardForkEnabled(true)

	suite.Run(t, &AnteTestSuite{
		AnteTestSuite: baseSuite,
	})

	// Re-run the tests with EIP-712 Legacy encodings to ensure backwards compatibility.
	// LegacyEIP712Extension should not be run with current TypedData encodings, since they are not compatible.
	suite.Run(t, &AnteTestSuite{
		AnteTestSuite:            baseSuite,
		useLegacyEIP712TypedData: true,
	})
}
