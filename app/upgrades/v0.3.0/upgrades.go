package v0_3_0

import (
	"errors"
	"fmt"
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		// Correct Pool Total stake
		pools := registryKeeper.GetAllPool(ctx)

		for _, pool := range pools {

			fmt.Printf("Fix Pool: %s\n", pool.Name)

			totalStake := uint64(0)
			totalFunds := uint64(0)

			for _, stakerAddress := range pool.Stakers {
				staker, found := registryKeeper.GetStaker(ctx, stakerAddress, pool.Id)
				if !found {
					panic("Error: Unexpected Error; Staker does not exist")
				}
				totalStake += staker.Amount
			}
			for _, funderAddress := range pool.Funders {
				funder, found := registryKeeper.GetFunder(ctx, funderAddress, pool.Id)
				if !found {
					panic("Error: Unexpected Error; Funder does not exist")
				}
				totalFunds += funder.Amount
			}

			pool.TotalStake = totalStake
			pool.TotalFunds = totalFunds

			registryKeeper.SetPool(ctx, pool)
		}

		// Unstake all users,
		delegators := registryKeeper.GetAllDelegator(ctx)

		for index, delegator := range delegators {
			// Withdraw all rewards for the sender.
			reward, err := customWithdraw(*registryKeeper, ctx, delegator.Id, delegator.Staker, delegator.Delegator)
			if err != nil {
				fmt.Printf("#%d: \tTransfer %d$KYVE to %s - ERROR : %s\n", index, reward, delegator.Delegator, err.Error())
				continue
			}

			// Transfer tokens from this module to sender.
			err = registryKeeper.TransferToAddress(ctx, delegator.Delegator, reward)
			if err != nil {
				return nil, err
			}

			fmt.Printf("#%d: \tTransfer %d$KYVE to %s\n", index, reward, delegator.Delegator)
		}

		for _, delegator := range registryKeeper.GetAllDelegator(ctx) {
			registryKeeper.RemoveDelegator(ctx, delegator.Id, delegator.Staker, delegator.Delegator)
		}

		for _, delegationPoolData := range registryKeeper.GetAllDelegationPoolData(ctx) {
			registryKeeper.RemoveDelegationPoolData(ctx, delegationPoolData.Id, delegationPoolData.Staker)
		}

		for _, delegationEntry := range registryKeeper.GetAllDelegationEntries(ctx) {
			registryKeeper.RemoveDelegationEntries(ctx, delegationEntry.Id, delegationEntry.Staker, delegationEntry.KIndex)
		}

		// Perform Re-Delegation
		for index, delegator := range delegators {
			f1Distribution := registrykeeper.CreateF1(*registryKeeper, ctx, delegator.Id, delegator.Staker, delegator.Delegator)
			// Withdraw all rewards for the sender.
			f1Distribution.Delegate(delegator.DelegationAmount)

			fmt.Printf("#%d: \tDelegate %d$KYVE from %s to %s\n", index, delegator.DelegationAmount, delegator.Delegator, delegator.Staker)
		}

		return vm, nil
	}
}

func customWithdraw(k registrykeeper.Keeper, ctx sdk.Context, poolId uint64, stakerAddress string, delegatorAddress string) (reward uint64, err error) {
	// Fetch metadata
	delegationPoolData, found := k.GetDelegationPoolData(ctx, poolId, stakerAddress)

	// Init default data-set, if this is the first delegator
	if !found {
		// Nothing to withdraw (CUSTOM ERROR HANDLING)
		return 0, errors.New("delegationPoolData not found")
	}

	delegator, found := k.GetDelegator(ctx, poolId, stakerAddress, delegatorAddress)
	if !found {
		// Nothing to withdraw (CUSTOM ERROR HANDLING)
		return 0, errors.New("delegator not found")
	}

	f1fMinus1, found := k.GetDelegationEntries(ctx, poolId, stakerAddress, delegationPoolData.LatestIndexK)
	f1fMinus1Balance := sdk.NewDec(0)
	if found {
		f1fMinus1Balance, _ = sdk.NewDecFromStr(f1fMinus1.Balance)
	}

	// F1Paper: T_f / n_f
	f1FinalBalance := sdk.NewDec(0)
	if delegationPoolData.TotalDelegation != 0 {
		decCurrentRewards := sdk.NewDec(int64(delegationPoolData.CurrentRewards))
		decTotalDelegation := sdk.NewDec(int64(delegationPoolData.TotalDelegation))

		f1FinalBalance = decCurrentRewards.Quo(decTotalDelegation)
	}

	// F1Paper Entry_f
	entryFBalance := f1fMinus1Balance.Add(f1FinalBalance)

	indexF := delegationPoolData.LatestIndexK + 1

	if delegationPoolData.LatestIndexWasUndelegation {
		//Remove Old entry
		k.RemoveDelegationEntries(ctx, poolId, stakerAddress, delegationPoolData.LatestIndexK)
	}

	// Insert Entry_F
	k.SetDelegationEntries(ctx, types.DelegationEntries{
		Id:      poolId,
		Balance: entryFBalance.String(),
		Staker:  stakerAddress,
		KIndex:  indexF,
	})

	delegationPoolData.LatestIndexWasUndelegation = false

	// Reset Values according to F1Paper, i.e T=0
	delegationPoolData.CurrentRewards = 0
	delegationPoolData.LatestIndexK = indexF

	k.SetDelegationPoolData(ctx, delegationPoolData)

	//Calculate Reward
	f1K, found := k.GetDelegationEntries(ctx, poolId, stakerAddress, delegator.KIndex)
	if !found {
		k.PanicHalt(ctx, "Delegator does not have entry")
	}

	//Remove Old entry
	k.RemoveDelegationEntries(ctx, poolId, stakerAddress, delegator.KIndex)

	//Update Delegator
	delegator.KIndex = indexF
	k.SetDelegator(ctx, delegator)

	f1kBalance, err := sdk.NewDecFromStr(f1K.Balance)
	if err != nil {
		return 0, err
	}

	return uint64((entryFBalance.Sub(f1kBalance)).Mul(sdk.NewDec(int64(delegator.DelegationAmount))).RoundInt64()), nil
}
