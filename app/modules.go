package app

import (
	// Cosmos SDK Utilities
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"

	// Bank
	"github.com/cosmos/cosmos-sdk/x/bank"

	// Capability
	"github.com/cosmos/cosmos-sdk/x/capability"

	// Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"

	// Distribution
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	// Evidence
	"github.com/cosmos/cosmos-sdk/x/evidence"

	// FeeGrant
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	// GenUtil
	"github.com/cosmos/cosmos-sdk/x/genutil"

	// Governance
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// IBC
	ibc "github.com/cosmos/ibc-go/v2/modules/core"

	// IBC Transfer
	"github.com/cosmos/ibc-go/v2/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"

	// Mint
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	// Parameters
	"github.com/cosmos/cosmos-sdk/x/params"

	// Registry
	"github.com/KYVENetwork/chain/x/registry"
	registrytypes "github.com/KYVENetwork/chain/x/registry/types"

	// Slashing
	"github.com/cosmos/cosmos-sdk/x/slashing"

	// Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

// appModuleBasics returns ModuleBasics for the module BasicManager.
var appModuleBasics = []module.AppModuleBasic{
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()...),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	registry.AppModuleBasic{},
	// this line is used by starport scaffolding # stargate/app/moduleBasic
}

// module account permissions
var moduleAccountPermissions = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	registrytypes.ModuleName:       {authtypes.Minter, authtypes.Burner, authtypes.Staking},
	// this line is used by starport scaffolding # stargate/app/maccPerms
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
