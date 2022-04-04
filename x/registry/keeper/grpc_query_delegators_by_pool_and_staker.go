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

func (k Keeper) AccountDelegationList(goCtx context.Context, req *types.QueryAccountDelegationListRequest) (*types.QueryAccountDelegationListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var delegated []types.DelegatorResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	delegatorStore := prefix.NewStore(store, types.KeyPrefix(types.DelegatorKeyPrefix))

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var delegator types.Delegator

		if err := k.cdc.Unmarshal(value, &delegator); err != nil {
			return false, nil
		}

		if delegator.Delegator != req.Address {
			return false, nil
		}

		if accumulate {

			pool, _ := k.GetPool(ctx, delegator.Id)

			f1 := F1Distribution{
				k:                k,
				ctx:              ctx,
				poolId:           pool.Id,
				stakerAddress:    delegator.Staker,
				delegatorAddress: delegator.Delegator,
			}

			delegationPoolData, _ := k.GetDelegationPoolData(ctx, pool.Id, delegator.Staker)

			delegated = append(delegated, types.DelegatorResponse{
				Account:            req.Address,
				Pool:               &pool,
				CurrentReward:      f1.getCurrentReward(),
				DelegationAmount:   delegator.DelegationAmount,
				Staker:             delegator.Staker,
				DelegationPoolData: &delegationPoolData,
			})
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountDelegationListResponse{
		Delegations: delegated,
		Pagination:  pageRes,
	}, nil
}
