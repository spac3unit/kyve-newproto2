package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetUnbondingDelegationQueueEntry ...
func (k Keeper) SetUnbondingDelegationQueueEntry(ctx sdk.Context, unbondingDelegationQueueEntry types.UnbondingDelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefix)
	b := k.cdc.MustMarshal(&unbondingDelegationQueueEntry)
	store.Set(types.UnbondingDelegationQueueEntryKey(
		unbondingDelegationQueueEntry.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefixIndex2)
	indexStore.Set(types.UnbondingDelegationQueueEntryKeyIndex2(
		unbondingDelegationQueueEntry.Delegator,
		unbondingDelegationQueueEntry.Index,
	), []byte{1})
}

// GetUnbondingDelegationQueueEntry returns a UnbondingDelegationQueueEntry from its index
func (k Keeper) GetUnbondingDelegationQueueEntry(ctx sdk.Context, index uint64) (val types.UnbondingDelegationQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefix)

	b := store.Get(types.UnbondingDelegationQueueEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingDelegationQueueEntry removes a UnbondingDelegationQueueEntry from the store
func (k Keeper) RemoveUnbondingDelegationQueueEntry(ctx sdk.Context, unbondingDelegationQueueEntry *types.UnbondingDelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefix)
	store.Delete(types.UnbondingDelegationQueueEntryKey(unbondingDelegationQueueEntry.Index))

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefixIndex2)
	indexStore.Delete(types.UnbondingDelegationQueueEntryKeyIndex2(
		unbondingDelegationQueueEntry.Delegator,
		unbondingDelegationQueueEntry.Index,
	))
}

// GetAllUnbondingDelegationQueueEntries returns all delegator unbondings
func (k Keeper) GetAllUnbondingDelegationQueueEntries(ctx sdk.Context) (list []types.UnbondingDelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UnbondingDelegationQueueEntryKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingDelegationQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ###################
// === QUEUE STATE ===
// ###################

// GetUnbondingDelegationQueueState returns the state for the undelegation queue
func (k Keeper) GetUnbondingDelegationQueueState(ctx sdk.Context) (state types.UnbondingDelegationQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := store.Get(types.UnbondingDelegationQueueStateKey)

	if b == nil {
		return state
	}

	k.cdc.MustUnmarshal(b, &state)
	return
}

// SetUnbondingDelegationQueueState saves the undelegation queue state
func (k Keeper) SetUnbondingDelegationQueueState(ctx sdk.Context, state types.UnbondingDelegationQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := k.cdc.MustMarshal(&state)
	store.Set(types.UnbondingDelegationQueueStateKey, b)
}
