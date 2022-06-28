package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetUnbondingStakingQueueEntry ...
func (k Keeper) SetUnbondingStakingQueueEntry(ctx sdk.Context, unbondingStakingQueueEntry types.UnbondingStakingQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefix)
	b := k.cdc.MustMarshal(&unbondingStakingQueueEntry)
	store.Set(types.UnbondingStakingQueueEntryKey(
		unbondingStakingQueueEntry.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefixIndex2)
	indexStore.Set(types.UnbondingStakingQueueEntryKeyIndex2(
		unbondingStakingQueueEntry.Staker,
		unbondingStakingQueueEntry.Index,
	), []byte{1})
}

// GetUnbondingStakingQueueEntry returns a UnbondingStakingQueueEntry from its index
func (k Keeper) GetUnbondingStakingQueueEntry(ctx sdk.Context, index uint64) (val types.UnbondingStakingQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefix)

	b := store.Get(types.UnbondingStakingQueueEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingStakingQueueEntry removes a UnbondingStakingQueueEntry from the store
func (k Keeper) RemoveUnbondingStakingQueueEntry(ctx sdk.Context, unbondingStakingQueueEntry *types.UnbondingStakingQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefix)
	store.Delete(types.UnbondingStakingQueueEntryKey(unbondingStakingQueueEntry.Index))

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefixIndex2)
	indexStore.Delete(types.UnbondingStakingQueueEntryKeyIndex2(
		unbondingStakingQueueEntry.Staker,
		unbondingStakingQueueEntry.Index,
	))
}

// GetAllUnbondingStakingQueueEntries returns all staker unbondings
func (k Keeper) GetAllUnbondingStakingQueueEntries(ctx sdk.Context) (list []types.UnbondingStakingQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakingQueueEntryKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingStakingQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ###################
// === QUEUE STATE ===
// ###################

// GetUnbondingStakingQueueState returns the state for the unstaking queue
func (k Keeper) GetUnbondingStakingQueueState(ctx sdk.Context) (state types.UnbondingStakingQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := store.Get(types.UnbondingStakingQueueStateKey)

	if b == nil {
		return state
	}

	k.cdc.MustUnmarshal(b, &state)
	return
}

// SetUnbondingStakingQueueState saves the unstaking queue state
func (k Keeper) SetUnbondingStakingQueueState(ctx sdk.Context, state types.UnbondingStakingQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := k.cdc.MustMarshal(&state)
	store.Set(types.UnbondingStakingQueueStateKey, b)
}

// ########################
// === UNBONDING STAKER ===
// ########################

// SetUnbondingStaker ...
func (k Keeper) SetUnbondingStaker(ctx sdk.Context, unbondingStaker types.UnbondingStaker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakerKeyPrefix)
	b := k.cdc.MustMarshal(&unbondingStaker)
	store.Set(types.UnbondingStakerKey(
		unbondingStaker.PoolId,
		unbondingStaker.Staker,
	), b)
}

// GetUnbondingStaker returns a UnbondingStaker from its address and pool
func (k Keeper) GetUnbondingStaker(ctx sdk.Context, poolId uint64, staker string) (val types.UnbondingStaker, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakerKeyPrefix)

	b := store.Get(types.UnbondingStakerKey(poolId, staker))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingStaker removes an UnbondingStaker from the store
func (k Keeper) RemoveUnbondingStaker(ctx sdk.Context, unbondingStaker *types.UnbondingStaker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakerKeyPrefix)
	store.Delete(types.UnbondingStakerKey(unbondingStaker.PoolId, unbondingStaker.Staker))
}

// GetAllUnbondingStakers returns all unbonding stakers
func (k Keeper) GetAllUnbondingStakers(ctx sdk.Context) (list []types.UnbondingStaker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingStakerKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingStaker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
