// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package testdata

import (
	contractutils "github.com/lyfeblocnetwork/os/contracts/utils"
	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
)

// LoadLBT9Contract load the LBT9 contract from the json representation of
// the Solidity contract.
func LoadLBT9Contract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("LBT9.json")
}
