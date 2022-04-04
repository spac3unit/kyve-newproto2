package keeper

import (
	"math"
	"math/rand"
	"sort"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// updateLowestFunder is an internal function that updates the lowest funder entry in a given pool.
func (k Keeper) updateLowestFunder(ctx sdk.Context, pool *types.Pool) {
	minAmount := uint64(math.Inf(0))
	minFunder := ""

	for _, account := range pool.Funders {
		funder, _ := k.GetFunder(ctx, account, pool.Id)

		if funder.Amount <= minAmount {
			minAmount = funder.Amount
			minFunder = funder.Account
		}
	}

	pool.LowestFunder = minFunder
}

// updateLowestStaker is an internal function that updates the lowest staker entry in a given pool.
func (k Keeper) updateLowestStaker(ctx sdk.Context, pool *types.Pool) {
	minAmount := uint64(math.Inf(0))
	minStaker := ""

	for _, account := range pool.Stakers {
		staker, _ := k.GetStaker(ctx, account, pool.Id)

		if staker.Amount <= minAmount {
			minAmount = staker.Amount
			minStaker = staker.Account
		}
	}

	pool.LowestStaker = minStaker
}

// removeFunder is an internal function that removes a funder from a given pool.
func (k Keeper) removeFunder(ctx sdk.Context, pool *types.Pool, funder *types.Funder) {
	// Find the index of the given funder.
	var funderIndex = -1

	for i, v := range pool.Funders {
		if v == funder.Account {
			funderIndex = i
			break
		}
	}

	// Return is the funder wasn't found.
	if funderIndex < 0 {
		return
	}

	// Remove funder from list of funders (replace with last entry and then slice).
	pool.Funders[funderIndex] = pool.Funders[len(pool.Funders)-1]
	pool.Funders = pool.Funders[:len(pool.Funders)-1]

	k.RemoveFunder(ctx, funder.Account, funder.PoolId)

	// Decrease the pool's total funds.
	pool.TotalFunds -= funder.Amount
}

// removeStaker is an internal function that removes a staker from a given pool.
func (k Keeper) removeStaker(ctx sdk.Context, pool *types.Pool, staker *types.Staker) {
	// Find the index of the given staker.
	var stakerIndex = -1

	for i, v := range pool.Stakers {
		if v == staker.Account {
			stakerIndex = i
			break
		}
	}

	// Return is the staker wasn't found.
	if stakerIndex < 0 {
		return
	}

	// Remove staker from list of stakers (replace with last entry and then slice).
	pool.Stakers[stakerIndex] = pool.Stakers[len(pool.Stakers)-1]
	pool.Stakers = pool.Stakers[:len(pool.Stakers)-1]

	k.RemoveStaker(ctx, staker.Account, staker.PoolId)

	// Decrease the pool's total stake.
	pool.TotalStake -= staker.Amount
}

// RandomChoiceCandidate ...
type RandomChoiceCandidate struct {
	Account string
	Amount  uint64
}

// getWeightedRandomChoice is an internal function that returns a random selection out of a list of candidates.
func (k Keeper) getWeightedRandomChoice(candidates []RandomChoiceCandidate, seed uint64) string {
	type WeightedRandomChoice struct {
		Elements    []string
		Weights     []uint64
		TotalWeight uint64
	}

	wrc := WeightedRandomChoice{}

	for _, candidate := range candidates {
		i := sort.Search(len(wrc.Weights), func(i int) bool { return wrc.Weights[i] > candidate.Amount })
		wrc.Weights = append(wrc.Weights, 0)
		wrc.Elements = append(wrc.Elements, "")
		copy(wrc.Weights[i+1:], wrc.Weights[i:])
		copy(wrc.Elements[i+1:], wrc.Elements[i:])
		wrc.Weights[i] = candidate.Amount
		wrc.Elements[i] = candidate.Account
		wrc.TotalWeight += candidate.Amount
	}

	rand.Seed(int64(seed))
	value := uint64(math.Floor(rand.Float64() * float64(wrc.TotalWeight)))

	for key, weight := range wrc.Weights {
		if weight > value {
			return wrc.Elements[key]
		}

		value -= weight
	}

	return ""
}

// Calculate Delegation weight to influnce the upload probability
// formula:
// A = 10000, dec = 10**9
// weight = dec * (sqrt(A * (A + x/dec)) - A)
func getDelegationWeight(delegation uint64) uint64 {

	const A uint64 = 10000

	number := A * (A + (delegation / 1_000_000_000))

	// Deterministic sqrt using only int
	// Uses the babylon recursive formula:
	// https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Babylonian_method
	var x uint64 = 14142 // expected value for 10000 $KYVE as input
	var xn uint64
	var epsilon uint64 = 100
	for epsilon > 2 {

		xn = (x + number/x) / 2

		if xn > x {
			epsilon = xn - x
		} else {
			epsilon = x - xn
		}
		x = xn
		println(x)
	}

	return (x - A) * 1_000_000_000
}

// getNextUploaderByRandom is an internal function that randomly selects the next uploader for a given pool.
func (k Keeper) getNextUploaderByRandom(ctx sdk.Context, pool *types.Pool) (nextUploader string) {
	var candidates []RandomChoiceCandidate

	for _, s := range pool.Stakers {
		staker, foundStaker := k.GetStaker(ctx, s, pool.Id)
		delegation, foundDelegation := k.GetDelegationPoolData(ctx, pool.Id, s)

		if foundStaker {
			if foundDelegation {
				candidates = append(candidates, RandomChoiceCandidate{
					Account: s,
					Amount:  staker.Amount + getDelegationWeight(delegation.TotalDelegation),
				})
			} else {
				candidates = append(candidates, RandomChoiceCandidate{
					Account: s,
					Amount:  staker.Amount,
				})
			}
		}
	}

	if len(candidates) == 0 {
		return pool.BundleProposal.NextUploader
	}

	return k.getWeightedRandomChoice(candidates, uint64(ctx.BlockHeight()+ctx.BlockTime().Unix()))
}

// slashStaker is an internal function that slashes a staker in a given pool by a certain percentage.
func (k Keeper) slashStaker(
	ctx sdk.Context, pool *types.Pool, stakerAddress string, slashAmountRatioDecimalString string,
) (slash uint64) {
	staker, found := k.GetStaker(ctx, stakerAddress, pool.Id)

	if found {
		// Parse the provided slash percentage and panic on any errors.
		slashAmountRatio, err := sdk.NewDecFromStr(slashAmountRatioDecimalString)
		if err != nil {
			panic("Invalid value for params: " + slashAmountRatioDecimalString + " error: " + err.Error())
		}

		// Compute how much we're going to slash the staker.
		slash = uint64(sdk.NewDec(int64(staker.Amount)).Mul(slashAmountRatio).RoundInt64())

		if staker.Amount == slash {
			// If we are slashing the entire staking amount, remove the staker.
			k.removeStaker(ctx, pool, &staker)
		} else {
			// Subtract slashing amount from staking amount, and update the pool's total stake.
			staker.Amount = staker.Amount - slash
			k.SetStaker(ctx, staker)

			pool.TotalStake -= slash
		}

		// Transfer the slashed amount to the treasury.
		err = k.transferToTreasury(ctx, slash)
		if err != nil {
			panic(err)
		}
	}

	return slash
}

// finalizeBundleProposal is an internal function that finalises the current bundle proposal.
func (k Keeper) finalizeBundleProposal(
	ctx sdk.Context, pool types.Pool, msgUploader string, msgBundleId string, msgByteSize uint64, msgBundleSize uint64,
) (*types.MsgSubmitBundleProposalResponse, error) {
	// Randomly select the next uploader (weighted random).
	nextUploader := k.getNextUploaderByRandom(ctx, &pool)

	// If the pool is in "genesis state" or an upload timeout has occurred, just update and return.
	if pool.BundleProposal.BundleId == "" {
		pool.BundleProposal = &types.BundleProposal{
			Uploader:     msgUploader,
			NextUploader: nextUploader,
			BundleId:     msgBundleId,
			ByteSize:     msgByteSize,
			FromHeight:   pool.BundleProposal.ToHeight,
			ToHeight:     pool.BundleProposal.ToHeight + msgBundleSize,
			CreatedAt:    uint64(ctx.BlockHeight()),
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// We have a currently active bundle, we now need to evaluate the round ...

	// Check if consensus has already been reached.
	validVotes := len(pool.BundleProposal.VotersValid)
	invalidVotes := len(pool.BundleProposal.VotersInvalid)

	valid := validVotes*2 > (len(pool.Stakers) - 1)
	invalid := invalidVotes*2 >= (len(pool.Stakers) - 1)

	if !valid && !invalid {
		// Not enough votes yet
		return nil, types.ErrQuorumNotReached
	}

	// Handle an empty bundle ...
	if pool.BundleProposal.BundleId == types.EmptyBundle {
		if valid {
			// Emit an empty bundle event.
			types.EmitEmptyBundleEvent(ctx, &pool)
		}

		if invalid {
			// Partially slash the uploader, because they were offline.
			slashAmount := k.slashStaker(ctx, &pool, pool.BundleProposal.Uploader, k.TimeoutSlash(ctx))

			// Emit a slash event for the uploader (TimeoutSlash).
			types.EmitSlashEvent(ctx, pool.Id, pool.BundleProposal.Uploader, slashAmount)

			// Update the current lowest staker.
			k.updateLowestStaker(ctx, &pool)

			// Emit a bundle timeout event.
			types.EmitBundleTimeoutEvent(ctx, &pool)
		}

		pool.BundleProposal = &types.BundleProposal{
			Uploader:     msgUploader,
			NextUploader: nextUploader,
			BundleId:     msgBundleId,
			ByteSize:     msgByteSize,
			FromHeight:   pool.BundleProposal.ToHeight,
			ToHeight:     pool.BundleProposal.ToHeight + msgBundleSize,
			CreatedAt:    uint64(ctx.BlockHeight()),
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// Handle a valid bundle ...
	if valid {
		// Calculate the total reward for the bundle, and individual payouts.
		bundleReward := pool.OperatingCost + (pool.BundleProposal.ByteSize * k.StorageCost(ctx))

		networkFee, err := sdk.NewDecFromStr(k.NetworkFee(ctx))
		if err != nil {
			panic("Invalid value for params: " + err.Error())
		}

		treasuryPayout := uint64(sdk.NewDec(int64(bundleReward)).Mul(networkFee).RoundInt64())
		uploaderPayout := bundleReward - treasuryPayout

		// Calculate the delegation rewards for the uploader.
		uploader, foundUploader := k.GetStaker(ctx, pool.BundleProposal.Uploader, pool.Id)
		uploaderDelegation, foundUploaderDelegation := k.GetDelegationPoolData(ctx, pool.Id, pool.BundleProposal.Uploader)

		if foundUploader && foundUploaderDelegation {
			// If the uploader has no delegators, it keeps the delegation reward.

			if uploaderDelegation.DelegatorCount > 0 {
				// Calculate the reward, factoring in the node commission, and subtract from the uploader payout.
				commission, _ := sdk.NewDecFromStr(uploader.Commission)
				delegationReward := uint64(
					sdk.NewDec(int64(uploaderPayout)).Mul(sdk.NewDec(1).Sub(commission)).RoundInt64(),
				)

				uploaderPayout -= delegationReward
				uploaderDelegation.CurrentRewards += delegationReward

				k.SetDelegationPoolData(ctx, uploaderDelegation)
			}
		}

		// Calculate the individual cost for each pool funder.
		// NOTE: Because of integer division, it is possible that there is a small remainder.
		//       This remainder is in worst case MaxFundersAmount(tkyve) and is charged to the lowest funder.
		fundersCost := bundleReward / uint64(len(pool.Funders))
		fundersCostRemainder := bundleReward - uint64(len(pool.Funders))*fundersCost

		// Fetch the lowest funder, and find a new one if the current one isn't found.
		lowestFunder, foundLowestFunder := k.GetFunder(ctx, pool.LowestFunder, pool.Id)

		if !foundLowestFunder {
			k.updateLowestFunder(ctx, &pool)
			lowestFunder, _ = k.GetFunder(ctx, pool.LowestFunder, pool.Id)
		}

		// Remove every funder who can't afford the funder cost.
		if fundersCost+fundersCostRemainder > lowestFunder.Amount {
			// First, let's remove the lowest funder.
			k.removeFunder(ctx, &pool, &lowestFunder)

			err := k.transferToTreasury(ctx, lowestFunder.Amount)
			if err != nil {
				return nil, err
			}

			// Now, let's remove all other funders who have run out of funds.
			for _, account := range pool.Funders {
				funder, _ := k.GetFunder(ctx, account, pool.Id)

				if funder.Amount < fundersCost {
					k.removeFunder(ctx, &pool, &funder)

					err := k.transferToTreasury(ctx, funder.Amount)
					if err != nil {
						return nil, err
					}
				}
			}

			// Recalculate the lowest funder, update, and return.
			k.updateLowestFunder(ctx, &pool)

			pool.BundleProposal = &types.BundleProposal{
				Uploader:      pool.BundleProposal.Uploader,
				NextUploader:  pool.BundleProposal.NextUploader,
				BundleId:      pool.BundleProposal.BundleId,
				ByteSize:      pool.BundleProposal.ByteSize,
				FromHeight:    pool.BundleProposal.FromHeight,
				ToHeight:      pool.BundleProposal.ToHeight,
				CreatedAt:     uint64(ctx.BlockHeight()),
				VotersValid:   pool.BundleProposal.VotersValid,
				VotersInvalid: pool.BundleProposal.VotersInvalid,
			}

			k.SetPool(ctx, pool)

			// Emit a bundle dropped event because of insufficient funds.
			types.EmitBundleDroppedInsufficientFundsEvent(ctx, &pool)

			return &types.MsgSubmitBundleProposalResponse{}, nil
		}

		// Charge every funder equally.
		for _, account := range pool.Funders {
			funder, _ := k.GetFunder(ctx, account, pool.Id)
			funder.Amount -= fundersCost
			k.SetFunder(ctx, funder)
		}

		// Remove any remainder cost from the lowest funder.
		lowestFunder, _ = k.GetFunder(ctx, pool.LowestFunder, pool.Id)
		lowestFunder.Amount -= fundersCostRemainder
		k.SetFunder(ctx, lowestFunder)

		// Subtract bundle reward from the pool's total funds.
		pool.TotalFunds -= bundleReward

		// Partially slash all nodes who voted incorrectly.
		for _, voter := range pool.BundleProposal.GetVotersInvalid() {
			slashAmount := k.slashStaker(ctx, &pool, voter, k.VoteSlash(ctx))
			types.EmitSlashEvent(ctx, pool.Id, voter, slashAmount)
		}

		// Send payout to treasury.
		errTreasury := k.transferToTreasury(ctx, treasuryPayout)
		if errTreasury != nil {
			return nil, errTreasury
		}

		// Send payout to uploader.
		errTransfer := k.transferToAddress(ctx, pool.BundleProposal.Uploader, uploaderPayout)
		if errTransfer != nil {
			return nil, errTransfer
		}

		// Finalise the proposal, saving useful information.
		pool.HeightArchived = pool.BundleProposal.ToHeight
		pool.BytesArchived = pool.BytesArchived + pool.BundleProposal.ByteSize
		pool.TotalBundleRewards = pool.TotalBundleRewards + bundleReward
		pool.TotalBundles = pool.TotalBundles + 1

		k.SetProposal(ctx, types.Proposal{
			BundleId:    pool.BundleProposal.BundleId,
			PoolId:      pool.Id,
			Uploader:    pool.BundleProposal.Uploader,
			FromHeight:  pool.BundleProposal.FromHeight,
			ToHeight:    pool.BundleProposal.ToHeight,
			FinalizedAt: uint64(ctx.BlockHeight()),
		})

		// Emit a valid bundle event.
		types.EmitBundleValidEvent(ctx, &pool, bundleReward)

		// Update and return.
		pool.BundleProposal = &types.BundleProposal{
			Uploader:     msgUploader,
			NextUploader: nextUploader,
			BundleId:     msgBundleId,
			ByteSize:     msgByteSize,
			FromHeight:   pool.BundleProposal.ToHeight,
			ToHeight:     pool.BundleProposal.ToHeight + msgBundleSize,
			CreatedAt:    uint64(ctx.BlockHeight()),
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// Handle an invalid bundle ...
	if invalid {
		// Partially slash all nodes who voted incorrectly.
		for _, voter := range pool.BundleProposal.GetVotersValid() {
			slashAmount := k.slashStaker(ctx, &pool, voter, k.VoteSlash(ctx))
			types.EmitSlashEvent(ctx, pool.Id, voter, slashAmount)
		}

		// Partially slash the uploader.
		slashAmount := k.slashStaker(ctx, &pool, pool.BundleProposal.Uploader, k.UploadSlash(ctx))
		types.EmitSlashEvent(ctx, pool.Id, pool.BundleProposal.Uploader, slashAmount)

		// Update the current lowest staker.
		k.updateLowestStaker(ctx, &pool)

		// Emit an invalid bundle event.
		types.EmitBundleInvalidEvent(ctx, &pool)

		// Update and return.
		pool.BundleProposal = &types.BundleProposal{
			NextUploader: pool.BundleProposal.NextUploader,
			FromHeight:   pool.BundleProposal.FromHeight,
			ToHeight:     pool.BundleProposal.FromHeight,
			CreatedAt:    uint64(ctx.BlockHeight()),
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	return &types.MsgSubmitBundleProposalResponse{}, nil
}
