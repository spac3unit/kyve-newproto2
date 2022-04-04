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

func (k Keeper) AccountStakersDelegationList(goCtx context.Context, req *types.QueryAccountStakersDelegationListRequest) (*types.QueryAccountStakersDelegationListResponse, error) {
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
	delegatorStore := prefix.NewStore(store, types.KeyPrefix(types.DelegatorKeyPrefix))

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var delegator types.Delegator

		if err := k.cdc.Unmarshal(value, &delegator); err != nil {
			return false, nil
		}

		if delegator.Staker != req.Staker {
			return false, nil
		}

		if accumulate {

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

	delegationPoolData, _ := k.GetDelegationPoolData(ctx, pool.Id, req.Staker)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountStakersDelegationListResponse{
		Delegators:         delegators,
		Pool:               &pool,
		DelegationPoolData: &delegationPoolData,
		Pagination:         pageRes,
	}, nil
}
