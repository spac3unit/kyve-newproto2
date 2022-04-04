package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanPropose(goCtx context.Context, req *types.QueryCanProposeRequest) (*types.QueryCanProposeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load pool
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Check if pool is paused
	if pool.Paused {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Pool is paused",
		}, nil
	}

	// Check if sender is a staker in pool
	_, isStaker := k.GetStaker(ctx, req.Proposer, req.PoolId)
	if !isStaker {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Proposer is no staker",
		}, nil
	}

	// Check quorum
	if pool.BundleProposal.GetBundleId() != "" {
		validVotes := len(pool.BundleProposal.GetVotersValid())
		invalidVotes := len(pool.BundleProposal.GetVotersInvalid())

		if validVotes == 0 && invalidVotes == 0 {
			return &types.QueryCanProposeResponse{
				Possible: false,
				Reason:   "Quorum not reached yet",
			}, nil
		}

		valid := validVotes*2 > (len(pool.GetStakers()) - 1)
		invalid := invalidVotes*2 >= (len(pool.GetStakers()) - 1)

		if !valid && !invalid {
			return &types.QueryCanProposeResponse{
				Possible: false,
				Reason:   "Quorum not reached yet",
			}, nil
		}
	}

	// check if designated uploader
	if pool.BundleProposal.GetNextUploader() != req.Proposer {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Not designated uploader",
		}, nil
	}

	// check if pool has funds
	if pool.TotalFunds == 0 {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Pool has run out of funds",
		}, nil
	}

	// check if uploader is the only staker
	if len(pool.Stakers) < 2 {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Uploader is the only node in pool",
		}, nil
	}

	return &types.QueryCanProposeResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
