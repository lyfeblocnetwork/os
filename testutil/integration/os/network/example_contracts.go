// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package network

import (
	testconstants "github.com/lyfeblocnetwork/os/testutil/constants"
)

// chainsLBTHex is an utility map used to retrieve the LBT contract
// address in hex format from the chain ID.
//
// TODO: refactor to define this in the example chain initialization and pass as function argument
var chainsLBTHex = map[string]string{
	testconstants.ExampleChainID: testconstants.LBTContractMainnet,
}

// GetLBTContractHex returns the hex format of address for the LBT contract
// given the chainID. If the chainID is not found, it defaults to the mainnet
// address.
func GetLBTContractHex(chainID string) string {
	address, found := chainsLBTHex[chainID]

	// default to mainnet address
	if !found {
		address = chainsLBTHex[testconstants.ExampleChainID]
	}

	return address
}
