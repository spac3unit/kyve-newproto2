package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleUploadTimeout is an end block hook that triggers an upload timeout for every pool (if applicable).
func (k Keeper) HandleUploadTimeout(goCtx context.Context) {
	// Unwrap context and fetch all pools.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pools := k.GetAllPool(ctx)

	// Iterate over all pools.
	for _, pool := range pools {
		// Remove next uploader immediately if not enough nodes are online
		if len(pool.Stakers) < 2 && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Remove next uploader immediately if pool has no funds
		if pool.TotalFunds == 0 && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Remove next uploader immediately if pool is paused
		if pool.Paused && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Skip if we haven't reached the upload timeout.
		if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval + k.UploadTimeout(ctx)) {
			continue
		}

		// We now know that the pool is active and the upload timeout has been reached.
		// Now we slash and remove the current next_uploader and select a new one.

		staker, foundStaker := k.GetStaker(ctx, pool.BundleProposal.NextUploader, pool.Id)

		// skip timeout slash if staker is not found
		if foundStaker {
			// slash next_uploader for not uploading in time
			slashAmount := k.slashStaker(ctx, &pool, staker.Account, k.TimeoutSlash(ctx))

			// emit slashing event
			types.EmitSlashEvent(ctx, pool.Id, staker.Account, slashAmount)

			staker, foundStaker = k.GetStaker(ctx, pool.BundleProposal.NextUploader, pool.Id)

			// check if next uploader is still there or already removed
			if foundStaker {
				// Transfer remaining stake to account.
				k.TransferToAddress(ctx, staker.Account, staker.Amount)

				// remove current next_uploader
				k.removeStaker(ctx, &pool, &staker)
			}

			// Update current lowest staker
			k.updateLowestStaker(ctx, &pool)
		}

		nextUploader := ""

		if len(pool.Stakers) > 0 {
			nextUploader = k.getNextUploaderByRandom(ctx, &pool, pool.Stakers)
		}

		// update bundle proposal
		pool.BundleProposal.NextUploader = nextUploader
		pool.BundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())

		k.SetPool(ctx, pool)
	}
}
