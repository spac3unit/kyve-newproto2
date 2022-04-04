package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetUnbondingEntries set a specific unbondingEntries in the store from its index
func (k Keeper) SetUnbondingEntries(ctx sdk.Context, unbondingEntries types.UnbondingEntries) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefix))
	b := k.cdc.MustMarshal(&unbondingEntries)
	store.Set(types.UnbondingEntriesKey(
		unbondingEntries.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	store2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefixByDelegator))
	b2 := k.cdc.MustMarshal(&unbondingEntries)
	store2.Set(types.UnbondingEntriesByDelegatorKey(
		unbondingEntries.Delegator,
		unbondingEntries.Index,
	), b2)
}

// GetUnbondingEntries returns a unbondingEntries from its index
func (k Keeper) GetUnbondingEntries(
	ctx sdk.Context,
	index uint64,

) (val types.UnbondingEntries, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefix))

	b := store.Get(types.UnbondingEntriesKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingEntries removes a unbondingEntries from the store
func (k Keeper) RemoveUnbondingEntries(
	ctx sdk.Context,
	unbondingEntry *types.UnbondingEntries,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefix))
	store.Delete(types.UnbondingEntriesKey(
		unbondingEntry.Index,
	))

	store2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefixByDelegator))
	store2.Delete(types.UnbondingEntriesByDelegatorKey(
		unbondingEntry.Delegator,
		unbondingEntry.Index,
	))
}

// GetAllUnbondingEntries returns all unbondingEntries
func (k Keeper) GetAllUnbondingEntries(ctx sdk.Context) (list []types.UnbondingEntries) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UnbondingEntries
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
