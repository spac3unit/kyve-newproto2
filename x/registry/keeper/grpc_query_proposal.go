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

// Deprecated: Proposals is deprecated as the return order is depending on the random bundle id
// Proposals return all bundles for a given pool ordered by bundle_id (which is random)
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

// Proposal returns the validated Proposal for a given bundle_id
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

// ProposalByHeight returns the proposal which contains the requested height of the datasource.
func (k Keeper) ProposalByHeight(goCtx context.Context, req *types.QueryProposalByHeightRequest) (*types.QueryProposalByHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex2}.AInt(req.PoolId)
	proposalIndexStore := prefix.NewStore(ctx.KVStore(k.storeKey), proposalPrefixBuilder.Key)
	proposalIndexIterator := proposalIndexStore.ReverseIterator(nil, types.KeyPrefixBuilder{}.AInt(req.Height+1).Key)

	defer proposalIndexIterator.Close()

	if proposalIndexIterator.Valid() {

		bundleId := string(proposalIndexIterator.Value())

		proposal, found := k.GetProposal(ctx, bundleId)
		if found {
			if proposal.FromHeight <= req.Height && proposal.ToHeight > req.Height {
				return &types.QueryProposalByHeightResponse{
					Proposal: proposal,
				}, nil
			}
		}
	}

	return nil, status.Error(codes.NotFound, "no bundle found")
}

// ProposalSinceFinalizedAt returns all proposals since a given finalizedAt height.
func (k Keeper) ProposalSinceFinalizedAt(goCtx context.Context, req *types.QueryProposalSinceFinalizedAtRequest) (*types.QueryProposalSinceFinalizedAtResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex3}.AInt(req.PoolId)
	proposalIndexStore := prefix.NewStore(ctx.KVStore(k.storeKey), proposalPrefixBuilder.Key)

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{}
	}

	if req.Pagination.Key == nil {
		// Find optimal key for query
		proposalIndexIterator := proposalIndexStore.Iterator(types.KeyPrefixBuilder{}.AInt(req.FinalizedAt).Key, nil)
		defer proposalIndexIterator.Close()

		if proposalIndexIterator.Valid() {
			req.Pagination.Key = proposalIndexIterator.Key()
		} else {
			return nil, status.Error(codes.NotFound, "no bundle found")
		}
	}

	var proposals []types.Proposal

	pageRes, err := query.FilteredPaginate(proposalIndexStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			bundleId := string(value)
			proposal, found := k.GetProposal(ctx, bundleId)
			if !found {
				return false, status.Error(codes.Internal, "bundleId should exist: "+bundleId)
			}
			proposals = append(proposals, proposal)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProposalSinceFinalizedAtResponse{
		Proposals:  proposals,
		Pagination: pageRes,
	}, nil
}
