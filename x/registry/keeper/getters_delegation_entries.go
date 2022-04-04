package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDelegationEntries set a specific delegationEntries in the store from its index
func (k Keeper) SetDelegationEntries(ctx sdk.Context, delegationEntries types.DelegationEntries) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationEntriesKeyPrefix))
	b := k.cdc.MustMarshal(&delegationEntries)
	store.Set(types.DelegationEntriesKey(
		delegationEntries.Id,
		delegationEntries.Staker,
		delegationEntries.KIndex,
	), b)
}

// GetDelegationEntries returns a delegationEntries from its index
func (k Keeper) GetDelegationEntries(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
	kIndex uint64,

) (val types.DelegationEntries, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationEntriesKeyPrefix))

	b := store.Get(types.DelegationEntriesKey(
		poolId,
		stakerAddress,
		kIndex,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDelegationEntries removes a delegationEntries from the store
func (k Keeper) RemoveDelegationEntries(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
	kIndex uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationEntriesKeyPrefix))
	store.Delete(types.DelegationEntriesKey(
		poolId,
		stakerAddress,
		kIndex,
	))
}

// GetAllDelegationEntries returns all delegationEntries
func (k Keeper) GetAllDelegationEntries(ctx sdk.Context) (list []types.DelegationEntries) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationEntriesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationEntries
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
