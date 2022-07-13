package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UndelegatePool handles the logic of an SDK message that allows undelegation from a protocol node in a specified pool.
func (k msgServer) UndelegatePool(
	goCtx context.Context, msg *types.MsgUndelegatePool,
) (*types.MsgUndelegatePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Perform undelegation
	if err := k.Undelegate(ctx, msg.Staker, msg.Id, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Create Unbonding queue entry
	if unbondingError := k.StartUnbondingDelegator(ctx, msg.Id, msg.Staker, msg.Creator, msg.Amount); unbondingError != nil {
		return nil, unbondingError
	}

	// Event an undelegate event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventUndelegatePool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Node:    msg.Staker,
		Amount:  msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgUndelegatePoolResponse{}, nil
}
