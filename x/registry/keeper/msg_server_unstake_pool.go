package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UnstakePool handles the logic of an SDK message that allows protocol nodes to unstake from a specified pool.
func (k msgServer) UnstakePool(
	goCtx context.Context, msg *types.MsgUnstakePool,
) (*types.MsgUnstakePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, foundPool := k.GetPool(ctx, msg.Id)
	if !foundPool {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	if err := k.StartUnbondingStaker(ctx, msg.Id, msg.Creator, msg.Amount); err != nil {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrUnstakeTooHigh.Error(), msg.Id)
	}

	return &types.MsgUnstakePoolResponse{}, nil
}
