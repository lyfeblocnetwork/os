// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package contracts

import (
	_ "embed"

	contractutils "github.com/lyfeblocnetwork/os/contracts/utils"
	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
)

var (
	// ERC20MinterBurnerDecimalsJSON are the compiled bytes of the ERC20MinterBurnerDecimalsContract
	//
	//go:embed solidity/ERC20MinterBurnerDecimals.json
	ERC20MinterBurnerDecimalsJSON []byte

	// ERC20MinterBurnerDecimalsContract is the compiled erc20 contract
	ERC20MinterBurnerDecimalsContract evmtypes.CompiledContract
)

func init() {
	var err error
	if ERC20MinterBurnerDecimalsContract, err = contractutils.ConvertHardhatBytesToCompiledContract(
		ERC20MinterBurnerDecimalsJSON,
	); err != nil {
		panic(err)
	}
}
