package registry

import (
	"github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	// Set pool count
	k.SetPoolCount(ctx, genState.PoolCount)
	// Set all the funder
	for _, elem := range genState.FunderList {
		k.SetFunder(ctx, elem)
	}
	// Set all the staker
	for _, elem := range genState.StakerList {
		k.SetStaker(ctx, elem)
	}
	// Set all the delegator
	for _, elem := range genState.DelegatorList {
		k.SetDelegator(ctx, elem)
	}
	// Set all the delegationPoolData
	for _, elem := range genState.DelegationPoolDataList {
		k.SetDelegationPoolData(ctx, elem)
	}
	// Set all the delegationEntries
	for _, elem := range genState.DelegationEntriesList {
		k.SetDelegationEntries(ctx, elem)
	}
	// Set all the proposal
	for _, elem := range genState.ProposalList {
		k.SetProposal(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PoolList = k.GetAllPool(ctx)
	genesis.PoolCount = k.GetPoolCount(ctx)
	genesis.FunderList = k.GetAllFunder(ctx)
	genesis.StakerList = k.GetAllStaker(ctx)
	genesis.DelegatorList = k.GetAllDelegator(ctx)
	genesis.DelegationPoolDataList = k.GetAllDelegationPoolData(ctx)
	genesis.DelegationEntriesList = k.GetAllDelegationEntries(ctx)
	genesis.ProposalList = k.GetAllProposal(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
