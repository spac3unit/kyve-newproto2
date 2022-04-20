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

// AccountStakedList returns all pools (with additional data) a user has staked into.
func (k Keeper) AccountStakedList(goCtx context.Context, req *types.QueryAccountStakedListRequest) (*types.QueryAccountStakedListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var staked []types.Staked
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	// Build prefix. Store is already indexed in an optimal way
	prefixBuilder := types.KeyPrefixBuilder{Key: types.KeyPrefix(types.StakerKeyPrefix)}.AString(req.Address).Key
	stakerStore := prefix.NewStore(store, prefixBuilder)

	pageRes, err := query.FilteredPaginate(stakerStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {

		if accumulate {

			var staker types.Staker
			if err := k.cdc.Unmarshal(value, &staker); err != nil {
				return false, err
			}

			pool, _ := k.GetPool(ctx, staker.PoolId)

			staked = append(staked, types.Staked{
				Staker:  staker.Account,
				PoolId:  staker.PoolId,
				Account: staker.Account,
				Amount:  staker.Amount,
				Pool:    &pool,
			})
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountStakedListResponse{
		Staked:     staked,
		Pagination: pageRes,
	}, nil
}
