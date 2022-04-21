package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UndelegatePool handles the logic of an SDK message that allows undelegation from a protocol node in a specified pool.
func (k msgServer) UndelegatePool(
	goCtx context.Context, msg *types.MsgUndelegatePool,
) (*types.MsgUndelegatePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is already a delegator.
	delegator, delegatorExists := k.GetDelegator(ctx, msg.Id, msg.Staker, msg.Creator)
	if !delegatorExists {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNotADelegator.Error(), msg.Id)
	}

	// Check if the sender is trying to undelegate more than they have delegated.
	if msg.Amount > delegator.DelegationAmount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrInsufficientFunds, types.ErrNotEnoughDelegation.Error(), msg.Amount)
	}

	// Create a new F1Distribution struct for interacting with delegations.
	f1Distribution := F1Distribution{
		k:                k.Keeper,
		ctx:              ctx,
		poolId:           msg.Id,
		stakerAddress:    msg.Staker,
		delegatorAddress: msg.Creator,
	}

	// Withdraw all rewards for the sender.
	reward := f1Distribution.Withdraw()

	// Transfer tokens from this module to sender.
	err := k.TransferToAddress(ctx, msg.Creator, reward)
	if err != nil {
		return nil, err
	}

	// Perform an internal re-delegation.
	undelegatedAmount := f1Distribution.Undelegate()
	redelegation := undelegatedAmount - msg.Amount
	f1Distribution.Delegate(redelegation)

	// Transfer the money
	err = k.TransferToAddress(ctx, msg.Creator, msg.Amount)
	if err != nil {
		k.PanicHalt(ctx, "Not enough money in module: "+err.Error())
	}

	// Event an undelegate event.
	types.EmitUndelegateEvent(ctx, pool.Id, delegator.Delegator, msg.Staker, msg.Amount)

	// Update and return.
	pool.TotalDelegation -= msg.Amount
	k.SetPool(ctx, pool)

	return &types.MsgUndelegatePoolResponse{}, nil
}
