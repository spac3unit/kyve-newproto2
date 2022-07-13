package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountRedelegation ...
func (k Keeper) AccountRedelegation(goCtx context.Context, req *types.QueryAccountRedelegationRequest) (*types.QueryAccountRedelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAccountRedelegationResponse{
		RedelegationCooldownEntries: k.GetRedelegationCooldownEntries(ctx, req.Address),
	}, nil
}
