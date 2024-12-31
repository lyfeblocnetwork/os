// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package constants

import (
	"fmt"

	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
)

const (
	// DefaultGasPrice is used in testing as the default to use for transactions
	DefaultGasPrice = 20

	// ExampleAttoDenom provides an example denom for use in tests
	ExampleAttoDenom = "alyfebloc"

	// ExampleMicroDenom provides an example denom for use in tests
	ExampleMicroDenom = "ulyfebloc"

	// ExampleDisplayDenom provides an example display denom for use in tests
	ExampleDisplayDenom = "lyfebloc"

	// ExampleBech32Prefix provides an example Bech32 prefix for use in tests
	ExampleBech32Prefix = "lyfebloc"

	// ExampleEIP155ChainID provides an example EIP-155 chain ID for use in tests
	ExampleEIP155ChainID = 9000

	// LBTContractMainnet is the LBT contract address for mainnet
	LBTContractMainnet = "0xD4949664cD82660AaE99bEdc034a0deA8A0bd517"
	// LBTContractTestnet is the LBT contract address for testnet
	LBTContractTestnet = "0xcc491f589b45d4a3c679016195b3fb87d7848210"
)

var (
	// ExampleChainIDPrefix provides a chain ID prefix for EIP-155 that can be used in tests
	ExampleChainIDPrefix = fmt.Sprintf("lyfebloc_%d", ExampleEIP155ChainID)

	// ExampleChainID provides a chain ID that can be used in tests
	ExampleChainID = ExampleChainIDPrefix + "-1"

	// SixDecimalsChainID provides a chain ID which is being set up with 6 decimals
	SixDecimalsChainID = "ossix_1776-2"

	// ExampleChainCoinInfo provides the coin info for the example chain
	//
	// It is a map of the chain id and its corresponding EvmCoinInfo
	// that allows initializing the app with different coin info based on the
	// chain id
	ExampleChainCoinInfo = map[string]evmtypes.EvmCoinInfo{
		ExampleChainID: {
			Denom:        ExampleAttoDenom,
			DisplayDenom: ExampleDisplayDenom,
			Decimals:     evmtypes.EighteenDecimals,
		},
		SixDecimalsChainID: {
			Denom:        ExampleMicroDenom,
			DisplayDenom: ExampleDisplayDenom,
			Decimals:     evmtypes.SixDecimals,
		},
	}
)
