package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Proposals(c context.Context, req *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var proposals []types.Proposal
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, types.KeyPrefix(types.ProposalKeyPrefix))

	pageRes, err := query.FilteredPaginate(proposalStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var proposal types.Proposal
		if err := k.cdc.Unmarshal(value, &proposal); err != nil {
			return false, err
		}

		// filter pool
		if proposal.PoolId != req.PoolId {
			return false, nil
		}

		if accumulate {
			proposals = append(proposals, proposal)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}

// Proposal returns validated Proposal for a given bundle_id
// This method is used to confirm if the provided bundle_id (e.g. from a third party)
// is an actual valid bundle uploaded to KYVE.
func (k Keeper) Proposal(c context.Context, req *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetProposal(
		ctx,
		req.BundleId,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryProposalResponse{Proposal: val}, nil
}
