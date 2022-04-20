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
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is the uploader, or the next uploader, in this pool.
	// TODO: Move to a custom error type.
	if msg.Creator == pool.BundleProposal.Uploader || msg.Creator == pool.BundleProposal.NextUploader {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, "Uploader and next uploader can't unstake.")
	}

	// Check if the sender has already voted on the current bundle.
	// TODO: Move to a custom error type.
	hasVotedValid, hasVotedInvalid := false, false

	for _, voter := range pool.BundleProposal.VotersValid {
		if msg.Creator == voter {
			hasVotedValid = true
		}
	}
	for _, voter := range pool.BundleProposal.VotersInvalid {
		if msg.Creator == voter {
			hasVotedInvalid = true
		}
	}

	if hasVotedValid || hasVotedInvalid {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, "Sender has already voted so can't unstake.")
	}

	// Check if the sender is a staker in this pool.
	staker, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.ErrNotFound
	}

	// Check if the sender is trying to unstake more than they have staked.
	if msg.Amount > staker.Amount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrUnstakeTooHigh.Error(), staker.Account)
	}

	// Update state variables (or completely remove if fully unstaking).
	if staker.Amount == msg.Amount {
		k.removeStaker(ctx, &pool, &staker)
	} else {
		staker.Amount -= msg.Amount
		k.SetStaker(ctx, staker)
	}

	// Transfer tokens from this module to sender.
	err := k.TransferToAddress(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Emit an unstake event.
	types.EmitUnstakeEvent(ctx, msg.Id, msg.Creator, msg.Amount)

	// Update and return.
	pool.TotalStake -= msg.Amount
	k.updateLowestStaker(ctx, &pool)
	k.SetPool(ctx, pool)

	return &types.MsgUnstakePoolResponse{}, nil
}
