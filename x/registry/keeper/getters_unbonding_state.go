package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetUnbondingState set unbondingState in the store
func (k Keeper) SetUnbondingState(ctx sdk.Context, unbondingState types.UnbondingState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingStateKey))
	b := k.cdc.MustMarshal(&unbondingState)
	store.Set([]byte{0}, b)
}

// GetUnbondingState returns unbondingState
func (k Keeper) GetUnbondingState(ctx sdk.Context) (val types.UnbondingState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUnbondingState removes unbondingState from the store
func (k Keeper) RemoveUnbondingState(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingStateKey))
	store.Delete([]byte{0})
}
