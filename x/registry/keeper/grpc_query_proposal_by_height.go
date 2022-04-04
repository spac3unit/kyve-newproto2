package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ProposalByHeight(goCtx context.Context, req *types.QueryProposalByHeightRequest) (*types.QueryProposalByHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposals := k.GetAllProposal(ctx)

	for _, proposal := range proposals {
		if proposal.PoolId == req.PoolId && proposal.FromHeight <= req.Height && proposal.ToHeight > req.Height {
			return &types.QueryProposalByHeightResponse{
				Proposal: proposal,
			}, nil
		}
	}

	return nil, status.Error(codes.NotFound, "no bundle found")
}
