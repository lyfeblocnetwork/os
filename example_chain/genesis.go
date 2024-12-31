// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package example_chain

import (
	"encoding/json"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	erc20types "github.com/lyfeblocnetwork/os/x/erc20/types"
	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
)

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by an identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewEVMGenesisState returns the default genesis state for the EVM module.
//
// NOTE: for the example chain implementation we need to set the default EVM denomination
// and enable ALL precompiles.
func NewEVMGenesisState() *evmtypes.GenesisState {
	evmGenState := evmtypes.DefaultGenesisState()
	evmGenState.Params.ActiveStaticPrecompiles = evmtypes.AvailableStaticPrecompiles

	return evmGenState
}

// NewErc20GenesisState returns the default genesis state for the ERC20 module.
//
// NOTE: for the example chain implementation we are also adding a default token pair,
// which is the base denomination of the chain (i.e. the LBT contract).
func NewErc20GenesisState() *erc20types.GenesisState {
	erc20GenState := erc20types.DefaultGenesisState()
	erc20GenState.TokenPairs = ExampleTokenPairs
	erc20GenState.Params.NativePrecompiles = append(erc20GenState.Params.NativePrecompiles, LBTContractMainnet)

	return erc20GenState
}

// NewMintGenesisState returns the default genesis state for the mint module.
//
// NOTE: for the example chain implementation we are also adding a default minter.
func NewMintGenesisState() *minttypes.GenesisState {
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params.MintDenom = ExampleChainDenom

	return mintGenState
}