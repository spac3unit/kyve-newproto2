package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// WithdrawPool handles the logic of an SDK message that allows delegators to collect all rewards from a specified pool.
func (k msgServer) WithdrawPool(
	goCtx context.Context, msg *types.MsgWithdrawPool,
) (*types.MsgWithdrawPoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is a delegator in this pool.
	_, isDelegator := k.GetDelegator(ctx, msg.Id, msg.Staker, msg.Creator)
	if !isDelegator {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNotADelegator.Error())
	}

	// Create a new F1Distribution struct for interacting with delegations.
	f1Distribution := F1Distribution{
		k:                k.Keeper,
		ctx:              ctx,
		poolId:           msg.Id,
		stakerAddress:    msg.Staker,
		delegatorAddress: msg.Creator,
	}

	// Withdraw all rewards for the sender.
	reward := f1Distribution.Withdraw()

	// Transfer tokens from this module to sender.
	if err := k.TransferToAddress(ctx, msg.Creator, reward); err != nil {
		return nil, err
	}

	return &types.MsgWithdrawPoolResponse{}, nil
}
