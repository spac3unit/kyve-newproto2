package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UnstakePool handles the logic of an SDK message that allows protocol nodes to unstake from a specified pool.
func (k msgServer) UnstakePool(goCtx context.Context, msg *types.MsgUnstakePool) (*types.MsgUnstakePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is a staker in this pool.
	staker, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.ErrNotFound
	}

	// Check if the sender is trying to unstake more than they have staked.
	// NOTE: Any amount that is still unbonding, isn't included.
	availableUnstakeAmount := uint64(0)
	if staker.UnbondingAmount <= staker.Amount {
		availableUnstakeAmount = staker.Amount - staker.UnbondingAmount
	}
	if msg.Amount > availableUnstakeAmount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrUnstakeTooHigh.Error(), staker.Account)
	}

	// Start the unbonding process.
	unbond := Unbond{k.Keeper, ctx}
	unbond.StartUnbond(pool.Id, msg.Creator, msg.Creator, msg.Amount)

	return &types.MsgUnstakePoolResponse{}, nil
}
