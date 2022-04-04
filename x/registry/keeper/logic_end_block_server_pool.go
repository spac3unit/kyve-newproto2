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
		// Skip if the pool is still in the "genesis state".
		if pool.BundleProposal.CreatedAt == 0 {
			continue
		}
		// Skip if the pool doesn't have any funds.
		if pool.TotalFunds == 0 {
			continue
		}
		// Skip if the pool is paused.
		if pool.Paused {
			continue
		}

		// Skip if we haven't reached the upload timeout.
		// TODO: Since timestamp is safe to use in Tendermint, consider switching.
		if (uint64(ctx.BlockHeight()) - pool.BundleProposal.CreatedAt) <= k.UploadTimeout(ctx) {
			continue
		}

		// We now know that the pool is active and the upload timeout has been reached.
		// Now we slash the current uploader and select a new one.

		// Fetch the votes on the current bundle and check if consensus has been reached.
		validVotes := len(pool.BundleProposal.VotersValid)
		invalidVotes := len(pool.BundleProposal.VotersInvalid)

		valid := validVotes*2 > (len(pool.Stakers) - 1)
		invalid := invalidVotes*2 >= (len(pool.Stakers) - 1)

		if valid || invalid {
			// If consensus was reached, let's finalise the current bundle.
			// TODO: What's the best way to handle this error? (since we're inside the EndBlock logic)
			k.finalizeBundleProposal(
				ctx, pool, pool.BundleProposal.NextUploader, types.EmptyBundle, 0, 0,
			)
		} else {
			// If consensus wasn't reached, we drop the bundle and emit an event.
			types.EmitBundleDroppedQuorumNotReachedEvent(ctx, &pool)

			pool.BundleProposal = &types.BundleProposal{
				NextUploader: pool.BundleProposal.NextUploader,
				FromHeight:   pool.BundleProposal.FromHeight,
				ToHeight:     pool.BundleProposal.FromHeight,
				CreatedAt:    uint64(ctx.BlockHeight()),
			}

			k.SetPool(ctx, pool)
		}
	}
}
