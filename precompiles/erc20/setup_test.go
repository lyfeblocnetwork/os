package erc20_test

import (
	"testing"

	erc20precompile "github.com/lyfeblocnetwork/os/precompiles/erc20"
	"github.com/lyfeblocnetwork/os/testutil/integration/os/factory"
	"github.com/lyfeblocnetwork/os/testutil/integration/os/grpc"
	testkeyring "github.com/lyfeblocnetwork/os/testutil/integration/os/keyring"
	"github.com/lyfeblocnetwork/os/testutil/integration/os/network"
	"github.com/stretchr/testify/suite"
)

var s *PrecompileTestSuite

// PrecompileTestSuite is the implementation of the TestSuite interface for ERC20 precompile
// unit tests.
type PrecompileTestSuite struct {
	suite.Suite

	bondDenom string
	// tokenDenom is the specific token denomination used in testing the ERC20 precompile.
	// This denomination is used to instantiate the precompile.
	tokenDenom  string
	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precompile *erc20precompile.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(PrecompileTestSuite)
	suite.Run(t, s)
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(2)
	integrationNetwork := network.NewUnitTestNetwork(
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	)
	grpcHandler := grpc.NewIntegrationHandler(integrationNetwork)
	txFactory := factory.New(integrationNetwork, grpcHandler)

	ctx := integrationNetwork.GetContext()
	sk := integrationNetwork.App.StakingKeeper
	bondDenom, err := sk.BondDenom(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(bondDenom, "bond denom cannot be empty")

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = integrationNetwork

	// Instantiate the precompile with an exemplary token denomination.
	//
	// NOTE: This has to be done AFTER assigning the suite fields.
	s.tokenDenom = "xmpl"
	s.precompile = s.setupERC20Precompile(s.tokenDenom)
}