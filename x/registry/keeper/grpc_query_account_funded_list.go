package keeper

import (
	"context"
	"fmt"
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountFundedList(goCtx context.Context, req *types.QueryAccountFundedListRequest) (*types.QueryAccountFundedListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var funded []types.Funded
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	funderStore := prefix.NewStore(store, types.KeyPrefix(types.FunderKeyPrefix))

	pageRes, err := query.FilteredPaginate(funderStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var funder types.Funder
		if err := k.cdc.Unmarshal(value, &funder); err != nil {
			return false, err
		}

		// filter account
		if funder.Account != req.Address {
			return false, nil
		}

		if accumulate {
			pool, _ := k.GetPool(ctx, funder.PoolId)

			funded = append(funded, types.Funded{
				// TODO Deprecated
				FundId:  fmt.Sprintf("%v %s", pool.Id, pool.LowestFunder),
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
