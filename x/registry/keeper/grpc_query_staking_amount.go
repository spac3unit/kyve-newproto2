package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakingAmount(goCtx context.Context, req *types.QueryStakingAmountRequest) (*types.QueryStakingAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	staker, exists := k.GetStaker(ctx, req.Staker, req.Id)
	if !exists {
		return nil, sdkErrors.ErrNotFound
	}

	return &types.QueryStakingAmountResponse{
		Amount: staker.Amount,
	}, nil
}
