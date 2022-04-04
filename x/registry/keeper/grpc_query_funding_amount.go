package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FundingAmount(goCtx context.Context, req *types.QueryFundingAmountRequest) (*types.QueryFundingAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	funder, exists := k.GetFunder(ctx, req.Funder, req.Id)
	if !exists {
		return nil, sdkErrors.ErrNotFound
	}

	return &types.QueryFundingAmountResponse{
		Amount: funder.Amount,
	}, nil
}
