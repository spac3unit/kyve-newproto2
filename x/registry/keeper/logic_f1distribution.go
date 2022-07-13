package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type F1Distribution struct {
	k                Keeper
	ctx              sdk.Context
	poolId           uint64
	stakerAddress    string
	delegatorAddress string
}

func (f1 F1Distribution) updateEntries(
	fMinus1Index uint64,
	currentRewards uint64,
	totalDelegation uint64,
	deleteOldEntry bool,
) (entryFBalance sdk.Dec, indexF uint64) {
	// F1Paper: Current period f = delegationPoolData.LatestIndexK + 1

	// get last but one entry for F1Distribution, init with zero if it is the first delegator
	// F1Paper: Entry_{f-1}
	f1fMinus1, found := f1.k.GetDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, fMinus1Index)
	f1fMinus1Balance := sdk.NewDec(0)
	if found {
		f1fMinus1Balance, _ = sdk.NewDecFromStr(f1fMinus1.Balance)
	}

	// F1Paper: T_f / n_f
	f1FinalBalance := sdk.NewDec(0)
	if totalDelegation != 0 {
		decCurrentRewards := sdk.NewDec(int64(currentRewards))
		decTotalDelegation := sdk.NewDec(int64(totalDelegation))

		f1FinalBalance = decCurrentRewards.Quo(decTotalDelegation)
	}

	// F1Paper Entry_f
	entryFBalance = f1fMinus1Balance.Add(f1FinalBalance)

	indexF = fMinus1Index + 1

	if deleteOldEntry {
		//Remove Old entry
		f1.k.RemoveDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, fMinus1Index)
	}

	// Insert Entry_F
	f1.k.SetDelegationEntries(f1.ctx, types.DelegationEntries{
		Id:      f1.poolId,
		Balance: entryFBalance.String(),
		Staker:  f1.stakerAddress,
		KIndex:  indexF,
	})

	return entryFBalance, indexF
}

func (f1 F1Distribution) Delegate(amount uint64) {

	if amount == 0 {
		return
	}

	// Fetch metadata
	delegationPoolData, found := f1.k.GetDelegationPoolData(f1.ctx, f1.poolId, f1.stakerAddress)

	// Init default data-set, if this is the first delegator
	if !found {
		delegationPoolData = types.DelegationPoolData{
			Id:              f1.poolId,
			CurrentRewards:  0,
			TotalDelegation: 0,
			LatestIndexK:    0,
			DelegatorCount:  0,
			Staker:          f1.stakerAddress,
		}
	}

	_, indexF := f1.updateEntries(delegationPoolData.LatestIndexK, delegationPoolData.CurrentRewards,
		delegationPoolData.TotalDelegation, delegationPoolData.LatestIndexWasUndelegation)

	delegationPoolData.LatestIndexWasUndelegation = false
	// Reset Values according to F1Paper, i.e T=0
	delegationPoolData.CurrentRewards = 0

	// Update metadata
	delegationPoolData.TotalDelegation += amount
	delegationPoolData.DelegatorCount += 1

	delegationPoolData.LatestIndexK = indexF

	f1.k.SetDelegationPoolData(f1.ctx, delegationPoolData)

	f1.k.SetDelegator(f1.ctx, types.Delegator{
		Staker:           f1.stakerAddress,
		Id:               f1.poolId,
		Delegator:        f1.delegatorAddress,
		DelegationAmount: amount,
		KIndex:           indexF,
	})
}

// Undelegate
// Undelegates the full amount.
// Withdraw() must be called before, otherwise the reward is gone
func (f1 F1Distribution) Undelegate() (undelegatedAmount uint64) {

	// Fetch metadata
	delegationPoolData, found := f1.k.GetDelegationPoolData(f1.ctx, f1.poolId, f1.stakerAddress)

	// Init default data-set, if this is the first delegator
	if !found {
		f1.k.PanicHalt(f1.ctx, "No delegationData although somebody is delegating")
	}

	delegator, found := f1.k.GetDelegator(f1.ctx, f1.poolId, f1.stakerAddress, f1.delegatorAddress)
	if !found {
		f1.k.PanicHalt(f1.ctx, "Not a delegator")
	}

	_, indexF := f1.updateEntries(delegationPoolData.LatestIndexK, delegationPoolData.CurrentRewards,
		delegationPoolData.TotalDelegation, delegationPoolData.LatestIndexWasUndelegation)

	// add flag that entry can be deleted after next entry is created
	delegationPoolData.LatestIndexWasUndelegation = true

	// Reset Values according to F1Paper, i.e T=0
	delegationPoolData.CurrentRewards = 0
	delegationPoolData.LatestIndexK = indexF

	// Update Metadata
	delegationPoolData.TotalDelegation -= delegator.DelegationAmount
	delegationPoolData.DelegatorCount -= 1

	//Remove Delegator
	f1.k.RemoveDelegator(f1.ctx, delegator.Id, delegator.Staker, delegator.Delegator)

	//Remove Old entry
	f1.k.RemoveDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, delegator.KIndex)

	if delegationPoolData.DelegatorCount == 0 {
		f1.k.RemoveDelegationPoolData(f1.ctx, delegationPoolData.Id, delegationPoolData.Staker)
		f1.k.RemoveDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, indexF)
	} else {
		f1.k.SetDelegationPoolData(f1.ctx, delegationPoolData)
	}

	return delegator.DelegationAmount
}

// Withdraw
// F1Withdraw updates the states for F1-Algorithm and returns the amount of coins the user has earned.
// The Method does NOT transfer the money.
func (f1 F1Distribution) Withdraw() (reward uint64) {
	// Fetch metadata
	delegationPoolData, found := f1.k.GetDelegationPoolData(f1.ctx, f1.poolId, f1.stakerAddress)

	// Init default data-set, if this is the first delegator
	if !found {
		f1.k.PanicHalt(f1.ctx, "No delegationData although somebody is delegating")
	}

	delegator, found := f1.k.GetDelegator(f1.ctx, f1.poolId, f1.stakerAddress, f1.delegatorAddress)
	if !found {
		f1.k.PanicHalt(f1.ctx, "Not a delegator")
	}

	entryFBalance, indexF := f1.updateEntries(delegationPoolData.LatestIndexK, delegationPoolData.CurrentRewards,
		delegationPoolData.TotalDelegation, delegationPoolData.LatestIndexWasUndelegation)

	delegationPoolData.LatestIndexWasUndelegation = false

	// Reset Values according to F1Paper, i.e T=0
	delegationPoolData.CurrentRewards = 0
	delegationPoolData.LatestIndexK = indexF

	f1.k.SetDelegationPoolData(f1.ctx, delegationPoolData)

	//Calculate Reward
	f1K, found := f1.k.GetDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, delegator.KIndex)
	if !found {
		f1.k.PanicHalt(f1.ctx, "Delegator does not have entry")
	}

	//Remove Old entry
	f1.k.RemoveDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, delegator.KIndex)

	//Update Delegator
	delegator.KIndex = indexF
	f1.k.SetDelegator(f1.ctx, delegator)

	f1kBalance, _ := sdk.NewDecFromStr(f1K.Balance)
	return uint64((entryFBalance.Sub(f1kBalance)).Mul(sdk.NewDec(int64(delegator.DelegationAmount))).RoundInt64())
}

// getCurrentReward
// Calculates and returns the current reward, *without* performing any state changes
func (f1 F1Distribution) getCurrentReward() (reward uint64) {

	delegator, found := f1.k.GetDelegator(f1.ctx, f1.poolId, f1.stakerAddress, f1.delegatorAddress)
	if !found {
		f1.k.PanicHalt(f1.ctx, "Not a delegator")
	}

	// Fetch metadata
	delegationPoolData, found := f1.k.GetDelegationPoolData(f1.ctx, f1.poolId, f1.stakerAddress)
	if !found {
		f1.k.PanicHalt(f1.ctx, "No delegationData although somebody is delegating")
	}

	// get last but one entry for F1Distribution, init with zero if it is the first delegator
	// F1Paper: Entry_{f-1}
	f1fMinus1, found := f1.k.GetDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, delegationPoolData.LatestIndexK)
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

	f1FinalBalance = f1FinalBalance.Add(f1fMinus1Balance)

	//Calculate Reward
	f1K, found := f1.k.GetDelegationEntries(f1.ctx, f1.poolId, f1.stakerAddress, delegator.KIndex)
	if !found {
		f1.k.PanicHalt(f1.ctx, "Delegator does not have entry")
	}

	f1kBalance, _ := sdk.NewDecFromStr(f1K.Balance)
	return uint64(sdk.NewDec(int64(delegator.DelegationAmount)).Mul(f1FinalBalance.Sub(f1kBalance)).RoundInt64())
}
