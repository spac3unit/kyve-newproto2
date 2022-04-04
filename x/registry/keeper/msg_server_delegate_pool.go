package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DelegatePool handles the logic of an SDK message that allows delegation to a protocol node from a specified pool.
func (k msgServer) DelegatePool(
	goCtx context.Context, msg *types.MsgDelegatePool,
) (*types.MsgDelegatePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is delegating to themselves.
	if msg.Creator == msg.Staker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrSelfDelegation.Error())
	}

	// Create a new F1Distribution struct for interacting with delegations.
	f1Distribution := F1Distribution{
		k:                k.Keeper,
		ctx:              ctx,
		poolId:           msg.Id,
		stakerAddress:    msg.Staker,
		delegatorAddress: msg.Creator,
	}

	// Check if the sender is already a delegator.
	_, delegatorExists := k.GetDelegator(ctx, msg.Id, msg.Staker, msg.Creator)

	if delegatorExists {
		// If the sender is already a delegator, first perform an undelegation, before then delegating.
		reward := f1Distribution.Withdraw()
		unDelegateAmount := f1Distribution.Undelegate()
		f1Distribution.Delegate(unDelegateAmount + msg.Amount)

		// Transfer tokens from sender to this module.
		amountToTransfer := unDelegateAmount + msg.Amount - reward

		err := k.transferToRegistry(ctx, msg.Creator, amountToTransfer)
		if err != nil {
			return nil, err
		}
	} else {
		// If the sender isn't already a delegator, simply create a new delegation entry.
		f1Distribution.Delegate(msg.Amount)

		// Transfer tokens from sender to this module.
		err := k.transferToRegistry(ctx, msg.Creator, msg.Amount)
		if err != nil {
			return nil, err
		}
	}

	// Emit a delegation event.
	types.EmitDelegateEvent(ctx, pool.Id, msg.Creator, msg.Staker, msg.Amount)

	// Update and return.
	pool.TotalDelegation += msg.Amount
	k.SetPool(ctx, pool)

	return &types.MsgDelegatePoolResponse{}, nil
}
