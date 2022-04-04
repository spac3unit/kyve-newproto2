package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Delegator(c context.Context, req *types.QueryDelegatorRequest) (*types.QueryDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetDelegator(
		ctx,
		req.PoolId,
		req.Staker,
		req.Delegator,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryDelegatorResponse{Delegator: val}, nil
}
