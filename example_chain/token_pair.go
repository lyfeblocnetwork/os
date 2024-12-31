// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package example_chain

import erc20types "github.com/lyfeblocnetwork/os/x/erc20/types"

// LBTContractMainnet is the LBT contract address for mainnet
const LBTContractMainnet = "0xD4949664cD82660AaE99bEdc034a0deA8A0bd517"

// ExampleTokenPairs creates a slice of token pairs, that contains a pair for the native denom of the example chain
// implementation.
var ExampleTokenPairs = []erc20types.TokenPair{
	{
		Erc20Address:  LBTContractMainnet,
		Denom:         ExampleChainDenom,
		Enabled:       true,
		ContractOwner: erc20types.OWNER_MODULE,
	},
}
