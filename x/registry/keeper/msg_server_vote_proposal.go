package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// VoteProposal handles the logic of an SDK message that allows protocol nodes to vote on a pool's bundle proposal.
func (k msgServer) VoteProposal(
	goCtx context.Context, msg *types.MsgVoteProposal,
) (*types.MsgVoteProposalResponse, error) {
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

	// Check if the sender is a protocol node (aka has staked into this pool).
	_, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Check if the sender is also the bundle's uploader.
	if pool.BundleProposal.GetUploader() == msg.Creator {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrVoterIsUploader.Error())
	}

	// Check if the sender is voting on a valid bundle.
	if msg.BundleId == "" || msg.BundleId != pool.BundleProposal.BundleId {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrNotFound, types.ErrInvalidBundleId.Error(), pool.BundleProposal.BundleId,
		)
	}

	// Check if the sender has already voted on the bundle.
	hasVotedValid, hasVotedInvalid := false, false

	for _, voter := range pool.BundleProposal.VotersValid {
		if voter == msg.Creator {
			hasVotedValid = true
		}
	}
	for _, voter := range pool.BundleProposal.VotersInvalid {
		if voter == msg.Creator {
			hasVotedInvalid = true
		}
	}

	if hasVotedValid || hasVotedInvalid {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrUnauthorized, types.ErrAlreadyVoted.Error(), pool.BundleProposal.BundleId,
		)
	}

	// Emit a vote event.
	types.EmitBundleVoteEvent(ctx, &pool, msg)

	// Update and return.
	if msg.Support {
		pool.BundleProposal.VotersValid = append(pool.BundleProposal.VotersValid, msg.Creator)
	} else {
		pool.BundleProposal.VotersInvalid = append(pool.BundleProposal.VotersInvalid, msg.Creator)
	}

	k.SetPool(ctx, pool)

	return &types.MsgVoteProposalResponse{}, nil
}
