package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDelegationPoolData set a specific delegationPoolData in the store from its index
func (k Keeper) SetDelegationPoolData(ctx sdk.Context, delegationPoolData types.DelegationPoolData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationPoolDataKeyPrefix))
	b := k.cdc.MustMarshal(&delegationPoolData)
	store.Set(types.DelegationPoolDataKey(
		delegationPoolData.Id,
		delegationPoolData.Staker,
	), b)
}

// GetDelegationPoolData returns a delegationPoolData from its index
func (k Keeper) GetDelegationPoolData(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
) (val types.DelegationPoolData, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationPoolDataKeyPrefix))

	b := store.Get(types.DelegationPoolDataKey(
		poolId,
		stakerAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDelegationPoolData removes a delegationPoolData from the store
func (k Keeper) RemoveDelegationPoolData(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationPoolDataKeyPrefix))
	store.Delete(types.DelegationPoolDataKey(
		poolId,
		stakerAddress,
	))
}

// GetAllDelegationPoolData returns all delegationPoolData
func (k Keeper) GetAllDelegationPoolData(ctx sdk.Context) (list []types.DelegationPoolData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegationPoolDataKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationPoolData
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
