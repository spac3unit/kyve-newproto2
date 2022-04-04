package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FundersList(goCtx context.Context, req *types.QueryFundersListRequest) (*types.QueryFundersListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var funders []*types.Funder

	// Load pool
	pool, found := k.GetPool(ctx, req.Id)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}

	for _, account := range pool.Funders {
		funder, _ := k.GetFunder(ctx, account, req.Id)
		funders = append(funders, &funder)
	}

	return &types.QueryFundersListResponse{
		Funders: funders,
	}, nil
}
