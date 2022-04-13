package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DelegatorsByPoolAndStaker returns all delegators for a specific pool that delegated to the given staker address
// It also returns the pool-object and the delegation information for that pool
// Pagination is supported
func (k Keeper) DelegatorsByPoolAndStaker(goCtx context.Context, req *types.QueryDelegatorsByPoolAndStakerRequest) (*types.QueryDelegatorsByPoolAndStakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	var delegators []types.StakerDelegatorResponse

	store := ctx.KVStore(k.storeKey)
	// Build prefix. Store is already indexed in an optimal way
	prefixBuilder := types.KeyPrefixBuilder{Key: types.KeyPrefix(types.DelegatorKeyPrefix)}.AInt(pool.Id).AString(req.Staker).Key
	delegatorStore := prefix.NewStore(store, prefixBuilder)

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {

		if accumulate {

			var delegator types.Delegator

			if err := k.cdc.Unmarshal(value, &delegator); err != nil {
				return false, nil
			}

			// Calculate current rewards for the delegator
			f1 := F1Distribution{
				k:                k,
				ctx:              ctx,
				poolId:           pool.Id,
				stakerAddress:    delegator.Staker,
				delegatorAddress: delegator.Delegator,
			}

			delegators = append(delegators, types.StakerDelegatorResponse{
				Delegator:        delegator.Delegator,
				CurrentReward:    f1.getCurrentReward(),
				DelegationAmount: delegator.DelegationAmount,
				Staker:           req.Staker,
			})
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationPoolData, _ := k.GetDelegationPoolData(ctx, pool.Id, req.Staker)

	return &types.QueryDelegatorsByPoolAndStakerResponse{
		Delegators:         delegators,
		Pool:               &pool,
		DelegationPoolData: &delegationPoolData,
		Pagination:         pageRes,
	}, nil
}
