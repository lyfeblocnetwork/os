// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)
package types

import (
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

func GetTxTypeName(txType int) string {
	switch txType {
	case gethtypes.DynamicFeeTxType:
		return "DynamicFeeTxType"
	case gethtypes.LegacyTxType:
		return "LegacyTxType"
	case gethtypes.AccessListTxType:
		return "AccessListTxType"
	default:
		panic("unknown tx type")
	}
}
