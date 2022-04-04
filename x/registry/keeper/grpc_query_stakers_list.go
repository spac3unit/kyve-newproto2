package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersList(goCtx context.Context, req *types.QueryStakersListRequest) (*types.QueryStakersListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var stakers []*types.StakerResponse

	// Load pool
	pool, found := k.GetPool(ctx, req.Id)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}

	for _, account := range pool.Stakers {
		staker, _ := k.GetStaker(ctx, account, req.Id)

		stakerResponse := types.StakerResponse{
			Staker:          staker.Account,
			PoolId:          staker.PoolId,
			Account:         staker.Account,
			Amount:          staker.Amount,
			UnbondingAmount: staker.UnbondingAmount,
			TotalDelegation: 0,
			Commission:      staker.Commission,
			Moniker:         staker.Moniker,
			Website:         staker.Website,
			Logo:            staker.Logo,
		}

		poolDelegationData, found := k.GetDelegationPoolData(ctx, staker.PoolId, staker.Account)

		if found {
			stakerResponse.TotalDelegation = poolDelegationData.TotalDelegation
		}

		stakers = append(stakers, &stakerResponse)
	}

	return &types.QueryStakersListResponse{
		Stakers: stakers,
	}, nil
}
