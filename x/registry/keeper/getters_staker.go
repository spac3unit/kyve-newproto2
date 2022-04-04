package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetStaker set a specific staker in the store from its index
func (k Keeper) SetStaker(ctx sdk.Context, staker types.Staker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	b := k.cdc.MustMarshal(&staker)
	store.Set(types.StakerKey(
		staker.Account,
		staker.PoolId,
	), b)
}

// GetStaker returns a staker from its index
func (k Keeper) GetStaker(
	ctx sdk.Context,
	staker string,
	poolId uint64,
) (val types.Staker, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))

	b := store.Get(types.StakerKey(
		staker,
		poolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveStaker removes a staker from the store
func (k Keeper) RemoveStaker(
	ctx sdk.Context,
	staker string,
	poolId uint64,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	store.Delete(types.StakerKey(
		staker,
		poolId,
	))
}

// GetAllStaker returns all staker
func (k Keeper) GetAllStaker(ctx sdk.Context) (list []types.Staker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
