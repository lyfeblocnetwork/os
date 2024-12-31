// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)
package ibc

import "errors"

var (
	ErrNoIBCVoucherDenom  = errors.New("denom is not an IBC voucher")
	ErrDenomTraceNotFound = errors.New("denom trace not found")
	ErrInvalidBaseDenom   = errors.New("invalid base denomination")
)
