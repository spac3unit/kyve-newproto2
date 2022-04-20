package v0_2_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
)

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Fetch the unbonding state.
		unbondingState, found := registryKeeper.GetUnbondingState(ctx)
		if !found {
			return vm, nil
		}

		// Check if queue is currently empty.
		if unbondingState.LowIndex == unbondingState.HighIndex {
			return vm, nil
		}

		// Iterate over queue and handle all entries.
		for i := unbondingState.LowIndex; i <= unbondingState.HighIndex; i++ {
			// Fetch the current entry.
			unbondingEntry, found := registryKeeper.GetUnbondingEntries(ctx, i)
			if !found {
				continue
			}

			// Ensure that we're only handling delegation unbondings.
			if unbondingEntry.Delegator != unbondingEntry.Staker {
				// Transfer tokens from registry to delegator.
				err := registryKeeper.TransferToAddress(ctx, unbondingEntry.Delegator, unbondingEntry.Amount)
				if err != nil {
					registryKeeper.PanicHalt(ctx, "Not enough money in registry module: "+err.Error())
				}
			}

			registryKeeper.RemoveUnbondingEntries(ctx, &unbondingEntry)
		}
		registryKeeper.RemoveUnbondingState(ctx)

		// Fix overflow in state of user kyve13xzj5e568n4kwe76ayzmzzuraz6c9vnaslrn3t and restore integrity

		// What happened?
		// During the unbonding-period the user got kicked out of the stakers-list
		// leading to the deletion of the stakers-entry. A short time after that, the user
		// rejoined the stakers list. Now the staker.UnbondingAmount was zero !!
		// This value then experienced an integer-overflow which then in a next step
		// allowed for the actual overflow in the staker amount.

		// The stake on block #153457 was 18_026_400_000_000 tkyve
		// On block #153458 the overflow occurred by unstaking 19_312_000_000_000 tkyve

		// Transfer original funds back ...
		registryKeeper.TransferToAddress(ctx, "kyve13xzj5e568n4kwe76ayzmzzuraz6c9vnaslrn3t", 18_026_400_000_000)
		// ... and delete staker entry and remove it from the pool
		registryKeeper.RemoveStaker(ctx, "kyve13xzj5e568n4kwe76ayzmzzuraz6c9vnaslrn3t", 1)

		pool, _ := registryKeeper.GetPool(ctx, 1)

		// Correct pool stake by adding the incorrect amount back and removing the correct amount
		pool.TotalStake = pool.TotalStake + uint64(19_312_000_000_000) - uint64(18_026_400_000_000)

		var stakerIndex = -1
		for i, v := range pool.Stakers {
			if v == "kyve13xzj5e568n4kwe76ayzmzzuraz6c9vnaslrn3t" {
				stakerIndex = i
				break
			}
		}

		// Remove staker from list of stakers (replace with last entry and then slice).
		pool.Stakers[stakerIndex] = pool.Stakers[len(pool.Stakers)-1]
		pool.Stakers = pool.Stakers[:len(pool.Stakers)-1]

		registryKeeper.SetPool(ctx, pool)

		// Return.
		return vm, nil
	}
}
