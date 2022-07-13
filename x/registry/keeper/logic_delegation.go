package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Delegate performs a safe delegation with all necessary checks
// Warning: does not transfer the amount (only the rewards)
func (k Keeper) Delegate(ctx sdk.Context, stakerAddress string, poolId uint64, delegatorAddress string, amount uint64) error {

	pool, found := k.GetPool(ctx, poolId)

	// Error if the pool isn't found.
	if !found {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
	}

	// Check if the sender is delegating to themselves.
	if delegatorAddress == stakerAddress {
		return sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrSelfDelegation.Error())
	}

	// Create a new F1Distribution struct for interacting with delegations.
	f1Distribution := F1Distribution{
		k:                k,
		ctx:              ctx,
		poolId:           poolId,
		stakerAddress:    stakerAddress,
		delegatorAddress: delegatorAddress,
	}

	// Check if the sender is already a delegator.
	_, delegatorExists := k.GetDelegator(ctx, poolId, stakerAddress, delegatorAddress)

	if delegatorExists {
		// If the sender is already a delegator, first perform an undelegation, before then delegating.
		reward := f1Distribution.Withdraw()
		err := k.TransferToAddress(ctx, delegatorAddress, reward)
		if err != nil {
			return err
		}

		// Perform redelegation
		unDelegateAmount := f1Distribution.Undelegate()
		f1Distribution.Delegate(unDelegateAmount + amount)

	} else {
		// If the sender isn't already a delegator, simply create a new delegation entry.
		f1Distribution.Delegate(amount)
	}

	// Update and return.
	pool.TotalDelegation += amount
	k.SetPool(ctx, pool)

	return nil
}

// Undelegate performs a safe undelegation
// Warning: It does not create an unbonding entry; it does not transfer the delegation back (only the rewards)
func (k Keeper) Undelegate(ctx sdk.Context, stakerAddress string, poolId uint64, delegatorAddress string, amount uint64) error {

	pool, poolFound := k.GetPool(ctx, poolId)
	if !poolFound {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
	}

	// Check if the sender is already a delegator.
	delegator, delegatorExists := k.GetDelegator(ctx, poolId, stakerAddress, delegatorAddress)
	if !delegatorExists {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNotADelegator.Error(), poolId)
	}

	// Check if the sender is trying to undelegate more than they have delegated.
	if amount > delegator.DelegationAmount {
		return sdkErrors.Wrapf(sdkErrors.ErrInsufficientFunds, types.ErrNotEnoughDelegation.Error(), amount)
	}

	// Create a new F1Distribution struct for interacting with delegations.
	f1Distribution := F1Distribution{
		k:                k,
		ctx:              ctx,
		poolId:           poolId,
		stakerAddress:    stakerAddress,
		delegatorAddress: delegatorAddress,
	}

	// Withdraw all rewards for the sender.
	reward := f1Distribution.Withdraw()

	// Transfer tokens from this module to sender.
	err := k.TransferToAddress(ctx, delegatorAddress, reward)
	if err != nil {
		return err
	}

	// Perform an internal re-delegation.
	undelegatedAmount := f1Distribution.Undelegate()
	redelegation := undelegatedAmount - amount
	f1Distribution.Delegate(redelegation)

	// Update and return.
	pool.TotalDelegation -= amount
	k.SetPool(ctx, pool)

	return nil
}
