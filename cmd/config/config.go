// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lyfeblocnetwork/os/types"
)

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(types.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                  // Shared
	config.SetFullFundraiserPath(types.BIP44HDPath) //nolint: staticcheck
}
