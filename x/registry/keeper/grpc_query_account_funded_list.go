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

// AccountFundedList returns all pools the given user has funded into.
// Supports Pagination.
func (k Keeper) AccountFundedList(goCtx context.Context, req *types.QueryAccountFundedListRequest) (*types.QueryAccountFundedListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var funded []types.Funded
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	// Build prefix. Store is already indexed in an optimal way
	prefixBuilder := types.KeyPrefixBuilder{Key: types.KeyPrefix(types.FunderKeyPrefix)}.AString(req.Address).Key
	funderStore := prefix.NewStore(store, prefixBuilder)

	pageRes, err := query.FilteredPaginate(funderStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {

		if accumulate {

			var funder types.Funder
			if err := k.cdc.Unmarshal(value, &funder); err != nil {
				return false, err
			}

			pool, _ := k.GetPool(ctx, funder.PoolId)

			funded = append(funded, types.Funded{
				Account: funder.Account,
				Amount:  funder.Amount,
				Pool:    &pool,
			})
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountFundedListResponse{
		Funded:     funded,
		Pagination: pageRes,
	}, nil
}
