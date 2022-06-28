package keeper

import (
	"context"
	"encoding/binary"
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountStakingUnbondings ...
func (k Keeper) AccountStakingUnbondings(goCtx context.Context, req *types.QueryAccountStakingUnbondingsRequest) (*types.QueryAccountStakingUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var stakingUnbondings []types.StakingUnbonding

	// Build prefix. Store is already indexed in an optimal way
	prefixBuilder := types.KeyPrefixBuilder{Key: types.UnbondingStakingQueueEntryKeyPrefixIndex2}.AString(req.Address).Key
	stakerUnbondingStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixBuilder)

	pageRes, err := query.FilteredPaginate(stakerUnbondingStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {

			index := binary.BigEndian.Uint64(key[0:8])
			unbondingEntry, _ := k.GetUnbondingStakingQueueEntry(ctx, index)

			pool, _ := k.GetPool(ctx, unbondingEntry.PoolId)

			stakingUnbondings = append(stakingUnbondings, types.StakingUnbonding{
				Amount:       unbondingEntry.Amount,
				CreationTime: unbondingEntry.CreationTime,
				Pool:         &pool,
			})
		}
		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountStakingUnbondingsResponse{
		Unbondings: stakingUnbondings,
		Pagination: pageRes,
	}, nil
}
