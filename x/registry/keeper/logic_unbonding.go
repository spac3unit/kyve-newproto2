package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Unbond struct {
	K   Keeper
	Ctx sdk.Context
}

// StartUnbond (Used for protocol undelegation and protocol unstaking)
// Creates an unbond-entry which will be executed in 21-days (and performing the real undelegation or unstaking)
// Also the amount is added as "unbondingAmount" to the staker or delegator to indicate the undelegation
// that is about to happen.
func (u Unbond) StartUnbond(poolId uint64, staker string, delegator string, amount uint64) {
	// unbondingState stores the start and the end of the queue with all unbonding entries
	// the queue is ordered by time
	unbondingState, found := u.K.GetUnbondingState(u.Ctx)
	if !found {
		unbondingState = types.UnbondingState{
			LowIndex:  0,
			HighIndex: 0,
		}
	}

	// Increase topIndex as a new entry is about to be appended
	newHighIndex := unbondingState.HighIndex + 1

	// UnbondingEntry stores all the information which are needed to perform
	// the undelegation at the end of the 21-days period
	unbondingEntry := types.UnbondingEntries{
		Index:        newHighIndex,
		PoolId:       poolId,
		Staker:       staker,
		Delegator:    delegator,
		CreationTime: uint64(u.Ctx.BlockTime().Unix()),
		Amount:       amount,
	}
	u.K.SetUnbondingEntries(u.Ctx, unbondingEntry)

	if delegator == staker {
		// If the delegator is also the staker (user delegated to himself)
		// the Staker-Store for the pool needs to be modified
		// the unbound amount is added to the UnbondingAmount to indicate that there are
		// pending unbondings.

		stakerEntry, _ := u.K.GetStaker(u.Ctx, staker, poolId)
		stakerEntry.UnbondingAmount += amount
		u.K.SetStaker(u.Ctx, stakerEntry)
	}

	// Update the unbonding state with the new head-position (highIndex)
	unbondingState.HighIndex = newHighIndex
	u.K.SetUnbondingState(u.Ctx, unbondingState)
}

// CheckAndPerformUndelegation is called at the end of every block and check the
// tail of the queue for Undelegations that can be performed
// This is usually O(1) or O(t) with t being the amount of undelegation-transactions which has been performed within
// a timeframe of one block (assuming KV-Store is O(1))
func (u Unbond) CheckAndPerformUndelegation() (count uint64, error error) {

	// Get Queue information
	unbondingState, found := u.K.GetUnbondingState(u.Ctx)
	if !found {
		return 0, nil
	}

	// Check if queue is currently empty
	if unbondingState.LowIndex == unbondingState.HighIndex {
		return 0, nil
	}

	// flag for computing every entry at the end of the queue which is due.
	undelegationPerformed := true
	// start processing the end of the queue
	for undelegationPerformed {
		undelegationPerformed = false

		// Get end of queue
		unbondingEntry, found := u.K.GetUnbondingEntries(u.Ctx, unbondingState.LowIndex+1)
		if !found {
			continue
		}

		// Check if unbonding time is over
		if unbondingEntry.CreationTime+uint64(types.UnbondingTime) < uint64(u.Ctx.BlockTime().Unix()) {

			if unbondingEntry.Staker == unbondingEntry.Delegator {
				// Handle case if delegator performed self-delegation

				stakerEntry, stakerExists := u.K.GetStaker(u.Ctx, unbondingEntry.Staker, unbondingEntry.PoolId)
				// if staker still exists perform the unstaking.
				// Otherwise, do nothing as he has already received his funds
				if stakerExists {

					pool, poolShouldExist := u.K.GetPool(u.Ctx, unbondingEntry.PoolId)
					if !poolShouldExist {
						u.K.PanicHalt(u.Ctx, "Pool should exist")
					}

					// Check if user got slashed during unbonding.
					// If user does now have less stake than in the unbonding entry
					// correct the UnbondingAmount to the stakers Amount
					if stakerEntry.UnbondingAmount > stakerEntry.Amount {
						unbondingEntry.Amount = stakerEntry.Amount
					}

					if stakerEntry.Amount == unbondingEntry.Amount {
						u.K.removeStaker(u.Ctx, &pool, &stakerEntry)
					} else {
						// Remove amount from current stake
						stakerEntry.Amount -= unbondingEntry.Amount

						// Remove amount from current unbondingAmount
						stakerEntry.UnbondingAmount -= unbondingEntry.Amount

						// Update Total stake in pool
						pool.TotalStake = pool.TotalStake - unbondingEntry.Amount

						// If staker still has stake -> update
						u.K.SetStaker(u.Ctx, stakerEntry)
					}

					// Update current lowest staker
					u.K.updateLowestStaker(u.Ctx, &pool)

					u.K.SetPool(u.Ctx, pool)

					// Transfer the money
					error := u.K.transferToAddress(u.Ctx, unbondingEntry.Delegator, unbondingEntry.Amount)
					if error != nil {
						u.K.PanicHalt(u.Ctx, "Not enough money in module: "+error.Error())
					}

					types.EmitUnstakeEvent(u.Ctx, pool.Id, unbondingEntry.Staker, unbondingEntry.Amount)
				}

			} else {
				// Handle case if delegator delegated to another staker

				// Transfer the money
				error := u.K.transferToAddress(u.Ctx, unbondingEntry.Delegator, unbondingEntry.Amount)
				if error != nil {
					u.K.PanicHalt(u.Ctx, "Not enough money in module: "+error.Error())
				}

			}

			// Remove entry from queue as it is no longer needed
			u.K.RemoveUnbondingEntries(u.Ctx, &unbondingEntry)

			// Update tailIndex (lowIndex) of queue
			unbondingState.LowIndex += 1
			u.K.SetUnbondingState(u.Ctx, unbondingState)

			// flags
			undelegationPerformed = true
			count += 1
		}
	}

	return count, nil
}
