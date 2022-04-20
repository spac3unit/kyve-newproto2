package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefundPool handles the logic of an SDK message that allows funders to defund from a specified pool.
func (k msgServer) DefundPool(goCtx context.Context, msg *types.MsgDefundPool) (*types.MsgDefundPoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is a funder in this pool.
	// TODO(@troy): Custom error?
	funder, isFunder := k.GetFunder(ctx, msg.Creator, msg.Id)
	if !isFunder {
		return nil, sdkErrors.ErrNotFound
	}

	// Check if the sender is trying to defund more than they have funded.
	if msg.Amount > funder.Amount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrDefundTooHigh.Error(), funder.Account)
	}

	// Update state variables (or completely remove if fully defunding).
	if funder.Amount == msg.Amount {
		k.removeFunder(ctx, &pool, &funder)
	} else {
		funder.Amount -= msg.Amount
		k.SetFunder(ctx, funder)
	}

	// Transfer tokens from this module to sender.
	err := k.TransferToAddress(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Emit a defund event.
	types.EmitDefundPoolEvent(ctx, msg.Id, msg.Creator, msg.Amount)

	// Update and return.
	pool.TotalFunds -= msg.Amount
	k.updateLowestFunder(ctx, &pool)
	k.SetPool(ctx, pool)

	return &types.MsgDefundPoolResponse{}, nil
}
