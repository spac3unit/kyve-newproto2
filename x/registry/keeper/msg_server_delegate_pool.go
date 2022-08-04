package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegatePool handles the logic of an SDK message that allows
// delegation to a protocol node from a specified pool.
func (k msgServer) DelegatePool(
	goCtx context.Context, msg *types.MsgDelegatePool,
) (*types.MsgDelegatePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Performs logical delegation without transferring the amount
	if err := k.Delegate(ctx, msg.Staker, msg.Id, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Transfer tokens from sender to this module.
	if transferErr := k.transferToRegistry(ctx, msg.Creator, msg.Amount); transferErr != nil {
		return nil, transferErr
	}

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventDelegatePool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Node:    msg.Staker,
		Amount:  msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgDelegatePoolResponse{}, nil
}
