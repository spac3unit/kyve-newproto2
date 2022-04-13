package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FundersList returns a list of all funders for given pool with their current funding amount
// This query is not paginated as it contains a maximum of types.MAX_FUNDERS entries
func (k Keeper) FundersList(goCtx context.Context, req *types.QueryFundersListRequest) (*types.QueryFundersListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryFundersListResponse{}

	// Load pool
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Fetch all funders information
	for _, account := range pool.Funders {
		funder, _ := k.GetFunder(ctx, account, req.PoolId)
		response.Funders = append(response.Funders, &funder)
	}

	return &response, nil
}
