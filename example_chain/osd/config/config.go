// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Bech32Prefix defines the Bech32 prefix used for accounts on the exemplary Lyfebloc blockchain.
	Bech32Prefix = "lyfebloc"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address.
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

const (
	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "lyfebloc"
	// BaseDenom defines to the default denomination used in the Lyfebloc example chain.
	BaseDenom = "alyfebloc"
	// BaseDenomUnit defines the precision of the base denomination.
	BaseDenomUnit = 18
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}