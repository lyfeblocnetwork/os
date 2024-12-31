package constants_test

import (
	"testing"

	"github.com/lyfeblocnetwork/os/testutil/constants"

	"github.com/lyfeblocnetwork/os/example_chain"
	chainconfig "github.com/lyfeblocnetwork/os/example_chain/osd/config"
	"github.com/stretchr/testify/require"
)

func TestRequireSameTestDenom(t *testing.T) {
	require.Equal(t,
		constants.ExampleAttoDenom,
		example_chain.ExampleChainDenom,
		"test denoms should be the same across the repo",
	)
}

func TestRequireSameTestBech32Prefix(t *testing.T) {
	require.Equal(t,
		constants.ExampleBech32Prefix,
		chainconfig.Bech32Prefix,
		"bech32 prefixes should be the same across the repo",
	)
}

func TestRequireSameLBTMainnet(t *testing.T) {
	require.Equal(t,
		constants.LBTContractMainnet,
		example_chain.LBTContractMainnet,
		"LBT contract addresses should be the same across the repo",
	)
}
