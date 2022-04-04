package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetFunder set a specific funder in the store from its index
func (k Keeper) SetFunder(ctx sdk.Context, funder types.Funder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))
	b := k.cdc.MustMarshal(&funder)
	store.Set(types.FunderKey(
		funder.Account,
		funder.PoolId,
	), b)
}

// GetFunder returns a funder from its index
func (k Keeper) GetFunder(
	ctx sdk.Context,
	Funder string,
	PoolId uint64,

) (val types.Funder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))

	b := store.Get(types.FunderKey(
		Funder,
		PoolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveFunder removes a funder from the store
func (k Keeper) RemoveFunder(
	ctx sdk.Context,
	Funder string,
	PoolId uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))
	store.Delete(types.FunderKey(
		Funder,
		PoolId,
	))
}

// GetAllFunder returns all funder
func (k Keeper) GetAllFunder(ctx sdk.Context) (list []types.Funder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Funder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
