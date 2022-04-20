package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferToAddress sends tokens from this module to a specified address.
func (k Keeper) TransferToAddress(ctx sdk.Context, address string, amount uint64) error {
	recipient, _ := sdk.AccAddressFromBech32(address)
	coins := sdk.NewCoins(sdk.NewInt64Coin("tkyve", int64(amount)))

	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	return err
}

// transferToRegistry sends tokens from a specified address to this module.
func (k Keeper) transferToRegistry(ctx sdk.Context, address string, amount uint64) error {
	sender, _ := sdk.AccAddressFromBech32(address)
	coins := sdk.NewCoins(sdk.NewInt64Coin("tkyve", int64(amount)))

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins)
	return err
}

// transferToTreasury sends tokens from this module to the treasury (community spend pool).
func (k Keeper) transferToTreasury(ctx sdk.Context, amount uint64) error {
	sender := k.accountKeeper.GetModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(sdk.NewInt64Coin("tkyve", int64(amount)))

	err := k.distrKeeper.FundCommunityPool(ctx, coins, sender)
	return err
}
