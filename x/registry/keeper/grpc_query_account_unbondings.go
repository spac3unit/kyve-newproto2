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

// AccountUnbondings returns all pending unbondings for the given account.
// They contain the information on when the unbonding occurs and who many $KYVE the user will receive
// Keep in mind that the unbonding amount for unstaking can be smaller when the user got slashed during the unbonding
func (k Keeper) AccountUnbondings(c context.Context, req *types.QueryAccountUnbondingsRequest) (*types.QueryAccountUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var unbondingEntries []types.UnbondingEntries
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	// Build prefix. Store is already indexed in an optimal way
	prefixBuilder := types.KeyPrefixBuilder{Key: types.KeyPrefix(types.UnbondingEntriesKeyPrefixByDelegator)}.AString(req.Address).Key
	unbondingEntriesStore := prefix.NewStore(store, prefixBuilder)

	pageRes, err := query.FilteredPaginate(unbondingEntriesStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {

		if accumulate {
			var unbondingEntry types.UnbondingEntries
			if err := k.cdc.Unmarshal(value, &unbondingEntry); err != nil {
				return false, err
			}

			unbondingEntries = append(unbondingEntries, unbondingEntry)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountUnbondingsResponse{UnbondingEntries: unbondingEntries, Pagination: pageRes}, nil
}
