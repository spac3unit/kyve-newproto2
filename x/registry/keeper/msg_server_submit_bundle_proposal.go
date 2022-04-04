package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubmitBundleProposal handles the logic of an SDK message that allows protocol nodes to submit a new bundle proposal.
func (k msgServer) SubmitBundleProposal(
	goCtx context.Context, msg *types.MsgSubmitBundleProposal,
) (*types.MsgSubmitBundleProposalResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}
	// Error if the pool is paused.
	if pool.Paused {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrPoolPaused.Error())
	}
	// Error if the pool has no funds.
	if len(pool.Funders) == 0 {
		return nil, sdkErrors.Wrap(sdkErrors.ErrInsufficientFunds, types.ErrFundsTooLow.Error())
	}

	// Check if the sender is a protocol node (aka has staked into this pool).
	_, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Check if the sender is the designated uploader.
	if pool.BundleProposal.NextUploader != msg.Creator {
		return nil, types.ErrNotDesignatedUploader
	}

	// Check to make sure enough nodes are online.
	if len(pool.Stakers) < 2 {
		return nil, types.ErrNotEnoughNodesOnline
	}

	// Validate bundle.
	if msg.BundleId == "" || msg.ByteSize == 0 {
		return nil, types.ErrInvalidArgs
	}
	if msg.BundleSize < pool.MinBundleSize {
		return nil, types.ErrMinBundleSize
	}

	// Call internal function that finalises the previous bundle and starts this one.
	return k.finalizeBundleProposal(ctx, pool, msg.Creator, msg.BundleId, msg.ByteSize, msg.BundleSize)
}
