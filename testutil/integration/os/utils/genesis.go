// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package utils

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	exampleapp "github.com/lyfeblocnetwork/os/example_chain"
	"github.com/lyfeblocnetwork/os/testutil/constants"
	testkeyring "github.com/lyfeblocnetwork/os/testutil/integration/os/keyring"
	"github.com/lyfeblocnetwork/os/testutil/integration/os/network"
	utiltx "github.com/lyfeblocnetwork/os/testutil/tx"
	erc20types "github.com/lyfeblocnetwork/os/x/erc20/types"
)

// CreateGenesisWithTokenPairs creates a genesis that includes
// the LBT and the provided denoms.
// If no denoms provided, creates only one dynamic precompile with the 'xmpl' denom.
func CreateGenesisWithTokenPairs(keyring testkeyring.Keyring, denoms ...string) network.CustomGenesisState {
	// Add all keys from the keyring to the genesis accounts as well.
	//
	// NOTE: This is necessary to enable the account to send EVM transactions,
	// because the Mono ante handler checks the account balance by querying the
	// account from the account keeper first. If these accounts are not in the genesis
	// state, the ante handler finds a zero balance because of the missing account.

	// if denom not provided, defaults to create only one dynamic erc20
	// precompile with the 'xmpl' denom
	if len(denoms) == 0 {
		denoms = []string{"xmpl"}
	}

	accs := keyring.GetAllAccAddrs()
	genesisAccounts := make([]*authtypes.BaseAccount, len(accs))
	for i, addr := range accs {
		genesisAccounts[i] = &authtypes.BaseAccount{
			Address:       addr.String(),
			PubKey:        nil,
			AccountNumber: uint64(i + 1), //nolint:gosec // G115
			Sequence:      1,
		}
	}

	accGenesisState := authtypes.DefaultGenesisState()
	for _, genesisAccount := range genesisAccounts {
		// NOTE: This type requires to be packed into a *types.Any as seen on SDK tests,
		// e.g. https://github.com/lyfebloc/cosmos-sdk/blob/v0.47.5-lyfebloc.2/x/auth/keeper/keeper_test.go#L193-L223
		accGenesisState.Accounts = append(accGenesisState.Accounts, codectypes.UnsafePackAny(genesisAccount))
	}

	// Add token pairs to genesis
	tokenPairs := make([]erc20types.TokenPair, 0, len(denoms)+1)
	tokenPairs = append(tokenPairs,
		// NOTE: the example token pairs are being added in the integration test utils
		exampleapp.ExampleTokenPairs...,
	)

	dynPrecAddr := make([]string, 0, len(denoms))
	for _, denom := range denoms {
		addr := utiltx.GenerateAddress().Hex()
		tp := erc20types.TokenPair{
			Erc20Address:  addr,
			Denom:         denom,
			Enabled:       true,
			ContractOwner: erc20types.OWNER_MODULE, // NOTE: Owner is the module account since it's a native token and was registered through governance
		}
		tokenPairs = append(tokenPairs, tp)
		dynPrecAddr = append(dynPrecAddr, addr)
	}

	// STR v2: update the NativePrecompiles and DynamicPrecompiles
	// with the LBT (default is mainnet) and 'xmpl' tokens in the erc20 params
	erc20GenesisState := exampleapp.NewErc20GenesisState()
	erc20GenesisState.TokenPairs = tokenPairs
	erc20GenesisState.Params.NativePrecompiles = []string{constants.LBTContractMainnet}
	erc20GenesisState.Params.DynamicPrecompiles = dynPrecAddr

	// Combine module genesis states
	return network.CustomGenesisState{
		authtypes.ModuleName:  accGenesisState,
		erc20types.ModuleName: erc20GenesisState,
	}
}