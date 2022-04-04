package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// FundPool handles the logic of an SDK message that allows funders to fund a specified pool.
func (k msgServer) FundPool(goCtx context.Context, msg *types.MsgFundPool) (*types.MsgFundPoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is already a funder.
	funder, funderExists := k.GetFunder(ctx, msg.Creator, msg.Id)

	if funderExists {
		funder.Amount += msg.Amount
		k.SetFunder(ctx, funder)
	} else {
		// Check if we have reached the maximum number of funders.
		// If we are funding more than the lowest funder, remove them.
		if len(pool.Funders) == types.MaxFunders {
			lowestFunder, _ := k.GetFunder(ctx, pool.LowestFunder, msg.Id)

			if msg.Amount > lowestFunder.Amount {
				// Transfer tokens from this module to the lowest funder.
				err := k.transferToAddress(ctx, lowestFunder.Account, lowestFunder.Amount)
				if err != nil {
					return nil, err
				}

				// Emit a defund event.
				types.EmitDefundPoolEvent(ctx, msg.Id, lowestFunder.Account, lowestFunder.Amount)

				// Remove lowest funder.
				k.removeFunder(ctx, &pool, &lowestFunder)
			} else {
				return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrFundsTooLow.Error(), lowestFunder.Amount)
			}
		}

		pool.Funders = append(pool.Funders, msg.Creator)
		k.SetFunder(ctx, types.Funder{
			Account: msg.Creator,
			PoolId:  msg.Id,
			Amount:  msg.Amount,
		})
	}

	// Transfer tokens from sender to this module.
	err := k.transferToRegistry(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Emit a fund event.
	types.EmitFundPoolEvent(ctx, msg)

	// Update and return.
	pool.TotalFunds += msg.Amount
	k.updateLowestFunder(ctx, &pool)
	k.SetPool(ctx, pool)

	return &types.MsgFundPoolResponse{}, nil
}
