// Copyright Lyfeloop Inc.(Lyfebloc)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/lyfeblocnetwork/lyfebloc/blob/main/LICENSE)

package network

import (
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	exampleapp "github.com/lyfeblocnetwork/os/example_chain"
	lyfebloctypes "github.com/lyfeblocnetwork/os/types"
	erc20types "github.com/lyfeblocnetwork/os/x/erc20/types"
	evmtypes "github.com/lyfeblocnetwork/os/x/evm/types"
	feemarkettypes "github.com/lyfeblocnetwork/os/x/feemarket/types"
)

// genSetupFn is the type for the module genesis setup functions
type genSetupFn func(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, customGenesis interface{}) (lyfebloctypes.GenesisState, error)

// defaultGenesisParams contains the params that are needed to
// setup the default genesis for the testing setup
type defaultGenesisParams struct {
	genAccounts []authtypes.GenesisAccount
	staking     StakingCustomGenesisState
	slashing    SlashingCustomGenesisState
	bank        BankCustomGenesisState
	gov         GovCustomGenesisState
	mint        MintCustomGenesisState
}

// genesisSetupFunctions contains the available genesis setup functions
// that can be used to customize the network genesis
var genesisSetupFunctions = map[string]genSetupFn{
	evmtypes.ModuleName:       genStateSetter[*evmtypes.GenesisState](evmtypes.ModuleName),
	erc20types.ModuleName:     genStateSetter[*erc20types.GenesisState](erc20types.ModuleName),
	govtypes.ModuleName:       genStateSetter[*govtypesv1.GenesisState](govtypes.ModuleName),
	feemarkettypes.ModuleName: genStateSetter[*feemarkettypes.GenesisState](feemarkettypes.ModuleName),
	distrtypes.ModuleName:     genStateSetter[*distrtypes.GenesisState](distrtypes.ModuleName),
	minttypes.ModuleName:      genStateSetter[*minttypes.GenesisState](minttypes.ModuleName),
	banktypes.ModuleName:      setBankGenesisState,
	authtypes.ModuleName:      setAuthGenesisState,
	consensustypes.ModuleName: func(_ *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, _ interface{}) (lyfebloctypes.GenesisState, error) {
		// no-op. Consensus does not have a genesis state on the application
		// but the params are used on it
		// (e.g. block max gas, max bytes).
		// This is handled accordingly on chain and context initialization
		return genesisState, nil
	},
	capabilitytypes.ModuleName: genStateSetter[*capabilitytypes.GenesisState](capabilitytypes.ModuleName),
}

// genStateSetter is a generic function to set module-specific genesis state
func genStateSetter[T proto.Message](moduleName string) genSetupFn {
	return func(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, customGenesis interface{}) (lyfebloctypes.GenesisState, error) {
		moduleGenesis, ok := customGenesis.(T)
		if !ok {
			return nil, fmt.Errorf("invalid type %T for %s module genesis state", customGenesis, moduleName)
		}

		genesisState[moduleName] = lyfeblocApp.AppCodec().MustMarshalJSON(moduleGenesis)
		return genesisState, nil
	}
}

// createValidatorSetAndSigners creates validator set with the amount of validators specified
// with the default power of 1.
func createValidatorSetAndSigners(numberOfValidators int) (*cmttypes.ValidatorSet, map[string]cmttypes.PrivValidator) {
	// Create validator set
	tmValidators := make([]*cmttypes.Validator, 0, numberOfValidators)
	signers := make(map[string]cmttypes.PrivValidator, numberOfValidators)

	for i := 0; i < numberOfValidators; i++ {
		privVal := mock.NewPV()
		pubKey, _ := privVal.GetPubKey()
		validator := cmttypes.NewValidator(pubKey, 1)
		tmValidators = append(tmValidators, validator)
		signers[pubKey.Address().String()] = privVal
	}

	return cmttypes.NewValidatorSet(tmValidators), signers
}

// createGenesisAccounts returns a slice of genesis accounts from the given
// account addresses.
func createGenesisAccounts(accounts []sdktypes.AccAddress) []authtypes.GenesisAccount {
	numberOfAccounts := len(accounts)
	genAccounts := make([]authtypes.GenesisAccount, 0, numberOfAccounts)
	for _, acc := range accounts {
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(
			acc, nil, 0, 0),
		)
	}
	return genAccounts
}

// getAccAddrsFromBalances returns a slice of genesis accounts from the
// given balances.
func getAccAddrsFromBalances(balances []banktypes.Balance) []sdktypes.AccAddress {
	numberOfBalances := len(balances)
	genAccounts := make([]sdktypes.AccAddress, 0, numberOfBalances)
	for _, balance := range balances {
		genAccounts = append(genAccounts, sdktypes.AccAddress(balance.Address))
	}
	return genAccounts
}

// createBalances creates balances for the given accounts and coin. Depending on
// the decimal representation of the denom, the amount is scaled to have the
// same value for all denoms.
func createBalances(
	accounts []sdktypes.AccAddress,
	denoms []string,
	denomsDecimals map[string]evmtypes.Decimals,
) []banktypes.Balance {
	numberOfAccounts := len(accounts)

	slices.Sort(denoms)

	coins := make([]sdktypes.Coin, len(denoms))
	for i, denom := range denoms {
		amount := PrefundedAccountInitialBalance
		dec, found := denomsDecimals[denom]
		// If the denom is not in the map, the 18 decimals representation is
		// used.
		if found {
			// amount is expressed in 18 decimals so it should be scaled down.
			amount = amount.Quo(dec.ConversionFactor())
		}
		coins[i] = sdktypes.NewCoin(denom, amount)
	}
	fundedAccountBalances := make([]banktypes.Balance, 0, numberOfAccounts)
	for _, acc := range accounts {
		balance := banktypes.Balance{
			Address: acc.String(),
			Coins:   coins,
		}

		fundedAccountBalances = append(fundedAccountBalances, balance)
	}
	return fundedAccountBalances
}

// createTestingApp creates an lyfebloc app
func createTestingApp(chainID string, customBaseAppOptions ...func(*baseapp.BaseApp)) *exampleapp.ExampleChain {
	// Create lyfebloc app
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	loadLatest := true
	appOptions := simutils.NewAppOptionsWithFlagHome(exampleapp.DefaultNodeHome)
	baseAppOptions := append(customBaseAppOptions, baseapp.SetChainID(chainID)) //nolint:gocritic

	return exampleapp.NewExampleApp(
		logger,
		db,
		nil,
		loadLatest,
		appOptions,
		baseAppOptions...,
	)
}

// createStakingValidator creates a staking validator from the given tm validator and bonded
func createStakingValidator(val *cmttypes.Validator, bondedAmt sdkmath.Int, operatorAddr *sdktypes.AccAddress) (stakingtypes.Validator, error) {
	pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey) //nolint:staticcheck
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	pkAny, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	opAddr := sdktypes.ValAddress(val.Address).String()
	if operatorAddr != nil {
		opAddr = sdktypes.ValAddress(operatorAddr.Bytes()).String()
	}

	// Default to 5% commission
	commission := stakingtypes.NewCommission(sdkmath.LegacyNewDecWithPrec(5, 2), sdkmath.LegacyNewDecWithPrec(2, 1), sdkmath.LegacyNewDecWithPrec(5, 2))
	validator := stakingtypes.Validator{
		OperatorAddress:   opAddr,
		ConsensusPubkey:   pkAny,
		Jailed:            false,
		Status:            stakingtypes.Bonded,
		Tokens:            bondedAmt,
		DelegatorShares:   sdkmath.LegacyOneDec(),
		Description:       stakingtypes.Description{},
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        commission,
		MinSelfDelegation: sdkmath.ZeroInt(),
	}
	return validator, nil
}

// createStakingValidators creates staking validators from the given tm validators and bonded
// amounts
func createStakingValidators(tmValidators []*cmttypes.Validator, bondedAmt sdkmath.Int, operatorsAddresses []sdktypes.AccAddress) ([]stakingtypes.Validator, error) {
	if len(operatorsAddresses) == 0 {
		return createStakingValidatorsWithRandomOperator(tmValidators, bondedAmt)
	}
	return createStakingValidatorsWithSpecificOperator(tmValidators, bondedAmt, operatorsAddresses)
}

// createStakingValidatorsWithRandomOperator creates staking validators with non-specified operator addresses.
func createStakingValidatorsWithRandomOperator(tmValidators []*cmttypes.Validator, bondedAmt sdkmath.Int) ([]stakingtypes.Validator, error) {
	amountOfValidators := len(tmValidators)
	stakingValidators := make([]stakingtypes.Validator, 0, amountOfValidators)
	for _, val := range tmValidators {
		validator, err := createStakingValidator(val, bondedAmt, nil)
		if err != nil {
			return nil, err
		}
		stakingValidators = append(stakingValidators, validator)
	}
	return stakingValidators, nil
}

// createStakingValidatorsWithSpecificOperator creates staking validators with the given operator addresses.
func createStakingValidatorsWithSpecificOperator(tmValidators []*cmttypes.Validator, bondedAmt sdkmath.Int, operatorsAddresses []sdktypes.AccAddress) ([]stakingtypes.Validator, error) {
	amountOfValidators := len(tmValidators)
	stakingValidators := make([]stakingtypes.Validator, 0, amountOfValidators)
	operatorsCount := len(operatorsAddresses)
	if operatorsCount != amountOfValidators {
		panic(fmt.Sprintf("provided %d validator operator keys but need %d!", operatorsCount, amountOfValidators))
	}
	for i, val := range tmValidators {
		validator, err := createStakingValidator(val, bondedAmt, &operatorsAddresses[i])
		if err != nil {
			return nil, err
		}
		stakingValidators = append(stakingValidators, validator)
	}
	return stakingValidators, nil
}

// createDelegations creates delegations for the given validators and account
func createDelegations(validators []stakingtypes.Validator, fromAccount sdktypes.AccAddress) []stakingtypes.Delegation {
	amountOfValidators := len(validators)
	delegations := make([]stakingtypes.Delegation, 0, amountOfValidators)
	for _, val := range validators {
		delegation := stakingtypes.NewDelegation(fromAccount.String(), val.OperatorAddress, sdkmath.LegacyOneDec())
		delegations = append(delegations, delegation)
	}
	return delegations
}

// getValidatorsSlashingGen creates the validators signingInfos and missedBlocks
// records necessary for the slashing module genesis
func getValidatorsSlashingGen(validators []stakingtypes.Validator, sk slashingtypes.StakingKeeper) (SlashingCustomGenesisState, error) {
	valCount := len(validators)
	signInfo := make([]slashingtypes.SigningInfo, valCount)
	missedBlocks := make([]slashingtypes.ValidatorMissedBlocks, valCount)
	for i, val := range validators {
		consAddrBz, err := val.GetConsAddr()
		if err != nil {
			return SlashingCustomGenesisState{}, err
		}
		consAddr, err := sk.ConsensusAddressCodec().BytesToString(consAddrBz)
		if err != nil {
			return SlashingCustomGenesisState{}, err
		}
		signInfo[i] = slashingtypes.SigningInfo{
			Address: consAddr,
			ValidatorSigningInfo: slashingtypes.ValidatorSigningInfo{
				Address: consAddr,
			},
		}
		missedBlocks[i] = slashingtypes.ValidatorMissedBlocks{
			Address: consAddr,
		}
	}
	return SlashingCustomGenesisState{
		signingInfo:  signInfo,
		missedBlocks: missedBlocks,
	}, nil
}

// StakingCustomGenesisState defines the staking genesis state
type StakingCustomGenesisState struct {
	denom string

	validators  []stakingtypes.Validator
	delegations []stakingtypes.Delegation
}

// setDefaultStakingGenesisState sets the default staking genesis state
func setDefaultStakingGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, overwriteParams StakingCustomGenesisState) lyfebloctypes.GenesisState {
	// Set staking params
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = overwriteParams.denom

	stakingGenesis := stakingtypes.NewGenesisState(
		stakingParams,
		overwriteParams.validators,
		overwriteParams.delegations,
	)
	genesisState[stakingtypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(stakingGenesis)
	return genesisState
}

type BankCustomGenesisState struct {
	totalSupply sdktypes.Coins
	balances    []banktypes.Balance
}

// setDefaultBankGenesisState sets the default bank genesis state
func setDefaultBankGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, overwriteParams BankCustomGenesisState) lyfebloctypes.GenesisState {
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		overwriteParams.balances,
		overwriteParams.totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	updatedBankGen := updateBankGenesisStateForChainID(*bankGenesis)
	genesisState[banktypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(&updatedBankGen)
	return genesisState
}

// SlashingCustomGenesisState defines the corresponding
// validators signing info and missed blocks for the genesis state
type SlashingCustomGenesisState struct {
	signingInfo  []slashingtypes.SigningInfo
	missedBlocks []slashingtypes.ValidatorMissedBlocks
}

// setDefaultSlashingGenesisState sets the default slashing genesis state
func setDefaultSlashingGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, overwriteParams SlashingCustomGenesisState) lyfebloctypes.GenesisState {
	slashingGen := slashingtypes.DefaultGenesisState()
	slashingGen.SigningInfos = overwriteParams.signingInfo
	slashingGen.MissedBlocks = overwriteParams.missedBlocks

	genesisState[slashingtypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(slashingGen)
	return genesisState
}

// setBankGenesisState updates the bank genesis state with custom genesis state
func setBankGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, customGenesis interface{}) (lyfebloctypes.GenesisState, error) {
	customGen, ok := customGenesis.(*banktypes.GenesisState)
	if !ok {
		return nil, fmt.Errorf("invalid type %T for bank module genesis state", customGenesis)
	}

	bankGen := &banktypes.GenesisState{}
	lyfeblocApp.AppCodec().MustUnmarshalJSON(genesisState[banktypes.ModuleName], bankGen)

	if len(customGen.Balances) > 0 {
		coins := sdktypes.NewCoins()
		bankGen.Balances = append(bankGen.Balances, customGen.Balances...)
		for _, b := range customGen.Balances {
			coins = append(coins, b.Coins...)
		}
		bankGen.Supply = bankGen.Supply.Add(coins...)
	}
	if len(customGen.DenomMetadata) > 0 {
		bankGen.DenomMetadata = append(bankGen.DenomMetadata, customGen.DenomMetadata...)
	}

	if len(customGen.SendEnabled) > 0 {
		bankGen.SendEnabled = append(bankGen.SendEnabled, customGen.SendEnabled...)
	}

	bankGen.Params = customGen.Params

	genesisState[banktypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(bankGen)
	return genesisState, nil
}

// calculateTotalSupply calculates the total supply from the given balances
func calculateTotalSupply(fundedAccountsBalances []banktypes.Balance) sdktypes.Coins {
	totalSupply := sdktypes.NewCoins()
	for _, balance := range fundedAccountsBalances {
		totalSupply = totalSupply.Add(balance.Coins...)
	}
	return totalSupply
}

// addBondedModuleAccountToFundedBalances adds bonded amount to bonded pool module account and include it on funded accounts
func addBondedModuleAccountToFundedBalances(
	fundedAccountsBalances []banktypes.Balance,
	totalBonded sdktypes.Coin,
) []banktypes.Balance {
	return append(fundedAccountsBalances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdktypes.Coins{totalBonded},
	})
}

// setDefaultAuthGenesisState sets the default auth genesis state
func setDefaultAuthGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, genAccs []authtypes.GenesisAccount) lyfebloctypes.GenesisState {
	defaultAuthGen := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(defaultAuthGen)
	return genesisState
}

// setAuthGenesisState updates the bank genesis state with custom genesis state
func setAuthGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, customGenesis interface{}) (lyfebloctypes.GenesisState, error) {
	customGen, ok := customGenesis.(*authtypes.GenesisState)
	if !ok {
		return nil, fmt.Errorf("invalid type %T for auth module genesis state", customGenesis)
	}

	authGen := &authtypes.GenesisState{}
	lyfeblocApp.AppCodec().MustUnmarshalJSON(genesisState[authtypes.ModuleName], authGen)

	if len(customGen.Accounts) > 0 {
		authGen.Accounts = append(authGen.Accounts, customGen.Accounts...)
	}

	authGen.Params = customGen.Params

	genesisState[authtypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(authGen)
	return genesisState, nil
}

// GovCustomGenesisState defines the gov genesis state
type GovCustomGenesisState struct {
	denom string
}

// setDefaultGovGenesisState sets the default gov genesis state
func setDefaultGovGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, overwriteParams GovCustomGenesisState) lyfebloctypes.GenesisState {
	govGen := govtypesv1.DefaultGenesisState()
	updatedParams := govGen.Params
	minDepositAmt := sdkmath.NewInt(1e18)
	updatedParams.MinDeposit = sdktypes.NewCoins(sdktypes.NewCoin(overwriteParams.denom, minDepositAmt))
	updatedParams.ExpeditedMinDeposit = sdktypes.NewCoins(sdktypes.NewCoin(overwriteParams.denom, minDepositAmt))
	govGen.Params = updatedParams
	genesisState[govtypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(govGen)
	return genesisState
}

// MintCustomGenesisState defines the gov genesis state
type MintCustomGenesisState struct {
	denom        string
	inflationMin sdkmath.LegacyDec
	inflationMax sdkmath.LegacyDec
}

// setDefaultGovGenesisState sets the default gov genesis state
//
// NOTE: for the testing network we don't want to have any minting
func setDefaultMintGenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState, overwriteParams MintCustomGenesisState) lyfebloctypes.GenesisState {
	mintGen := minttypes.DefaultGenesisState()
	updatedParams := mintGen.Params
	updatedParams.MintDenom = overwriteParams.denom
	updatedParams.InflationMin = overwriteParams.inflationMin
	updatedParams.InflationMax = overwriteParams.inflationMax

	mintGen.Params = updatedParams
	genesisState[minttypes.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(mintGen)
	return genesisState
}

func setDefaultErc20GenesisState(lyfeblocApp *exampleapp.ExampleChain, genesisState lyfebloctypes.GenesisState) lyfebloctypes.GenesisState {
	// NOTE: here we are using the setup from the example chain
	erc20Gen := exampleapp.NewErc20GenesisState()
	updatedErc20Gen := updateErc20GenesisStateForChainID(lyfeblocApp.ChainID(), *erc20Gen)

	genesisState[erc20types.ModuleName] = lyfeblocApp.AppCodec().MustMarshalJSON(&updatedErc20Gen)
	return genesisState
}

// defaultAuthGenesisState sets the default genesis state
// for the testing setup
func newDefaultGenesisState(lyfeblocApp *exampleapp.ExampleChain, params defaultGenesisParams) lyfebloctypes.GenesisState {
	genesisState := lyfeblocApp.DefaultGenesis()

	genesisState = setDefaultAuthGenesisState(lyfeblocApp, genesisState, params.genAccounts)
	genesisState = setDefaultStakingGenesisState(lyfeblocApp, genesisState, params.staking)
	genesisState = setDefaultBankGenesisState(lyfeblocApp, genesisState, params.bank)
	genesisState = setDefaultGovGenesisState(lyfeblocApp, genesisState, params.gov)
	genesisState = setDefaultSlashingGenesisState(lyfeblocApp, genesisState, params.slashing)
	genesisState = setDefaultMintGenesisState(lyfeblocApp, genesisState, params.mint)
	genesisState = setDefaultErc20GenesisState(lyfeblocApp, genesisState)

	return genesisState
}

// customizeGenesis modifies genesis state if there are any custom genesis state
// for specific modules
func customizeGenesis(lyfeblocApp *exampleapp.ExampleChain, customGen CustomGenesisState, genesisState lyfebloctypes.GenesisState) (lyfebloctypes.GenesisState, error) {
	var err error
	for mod, modGenState := range customGen {
		if fn, found := genesisSetupFunctions[mod]; found {
			genesisState, err = fn(lyfeblocApp, genesisState, modGenState)
			if err != nil {
				return genesisState, err
			}
		} else {
			panic(fmt.Sprintf("module %s not found in genesis setup functions", mod))
		}
	}
	return genesisState, err
}