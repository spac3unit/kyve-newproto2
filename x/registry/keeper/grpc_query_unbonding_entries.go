package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountUnbondings(c context.Context, req *types.QueryAccountUnbondingsRequest) (*types.QueryAccountUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var unbondingEntries []types.UnbondingEntries
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	unbondingEntriesStore := prefix.NewStore(store, append(types.KeyPrefix(types.UnbondingEntriesKeyPrefixByDelegator), types.KeyPrefix(req.Address)...))

	pageRes, err := query.FilteredPaginate(unbondingEntriesStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var unbondingEntry types.UnbondingEntries
		if err := k.cdc.Unmarshal(value, &unbondingEntry); err != nil {
			return false, err
		}

		if accumulate {
			unbondingEntries = append(unbondingEntries, unbondingEntry)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountUnbondingsResponse{UnbondingEntries: unbondingEntries, Pagination: pageRes}, nil
}
