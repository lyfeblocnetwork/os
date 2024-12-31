// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package grpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ethereum/go-ethereum/common"
	commongrpc "github.com/lyfeblocnetwork/os/testutil/integration/common/grpc"
	"github.com/lyfeblocnetwork/os/testutil/integration/os/network"
	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
	feemarkettypes "github.com/lyfeblocnetwork/os/x/feemarket/types"
)

// Handler is an interface that defines the methods that are used to query
// the network's modules via gRPC.
type Handler interface {
	commongrpc.Handler

	// EVM methods
	GetEvmAccount(address common.Address) (*evmtypes.QueryAccountResponse, error)
	EstimateGas(args []byte, GasCap uint64) (*evmtypes.EstimateGasResponse, error)
	GetEvmParams() (*evmtypes.QueryParamsResponse, error)
	GetEvmBaseFee() (*evmtypes.QueryBaseFeeResponse, error)
	GetBalanceFromEVM(address sdk.AccAddress) (*evmtypes.QueryBalanceResponse, error)

	// FeeMarket methods
	GetBaseFee() (*feemarkettypes.QueryBaseFeeResponse, error)
	GetFeeMarketParams() (*feemarkettypes.QueryParamsResponse, error)

	// Gov methods
	GetProposal(proposalID uint64) (*govtypes.QueryProposalResponse, error)
	GetGovParams(paramsType string) (*govtypes.QueryParamsResponse, error)
}

var _ Handler = (*IntegrationHandler)(nil)

// IntegrationHandler is a helper struct to query the network's modules
// via gRPC. This is to simulate the behavior of a real user and avoid querying
// the modules directly.
type IntegrationHandler struct {
	// We take the IntegrationHandler from common/grpc to get the common methods.
	*commongrpc.IntegrationHandler
	network network.Network
}

// NewIntegrationHandler creates a new IntegrationHandler instance.
func NewIntegrationHandler(network network.Network) Handler {
	return &IntegrationHandler{
		// Is there a better way to do this?
		IntegrationHandler: commongrpc.NewIntegrationHandler(network),
		network:            network,
	}
}