package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetProposal set a specific proposal in the store from its index
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	b := k.cdc.MustMarshal(&proposal)
	store.Set(types.ProposalKey(
		proposal.BundleId,
	), b)

	// Insert bundle id for second index
	storeIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	storeIndex.Set(types.ProposalKeyIndex2(proposal.PoolId, proposal.FromHeight), []byte(proposal.BundleId))

	// Insert bundle id for second index
	storeIndex3 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex3)
	storeIndex3.Set(types.ProposalKeyIndex3(proposal.PoolId, proposal.FinalizedAt), []byte(proposal.BundleId))
}

// GetProposal returns a proposal from its index
func (k Keeper) GetProposal(ctx sdk.Context, bundleId string) (val types.Proposal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))

	b := store.Get(types.ProposalKey(
		bundleId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveProposal removes a proposal from the store
func (k Keeper) RemoveProposal(ctx sdk.Context, proposal types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	store.Delete(types.ProposalKey(proposal.BundleId))

	indexStore2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	indexStore2.Delete(types.ProposalKeyIndex2(proposal.PoolId, proposal.FromHeight))

	// Insert bundle id for second index
	storeIndex3 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex3)
	storeIndex3.Delete(types.ProposalKeyIndex3(proposal.PoolId, proposal.FinalizedAt))
}

// GetAllProposal returns all proposal
func (k Keeper) GetAllProposal(ctx sdk.Context) (list []types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Proposal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
