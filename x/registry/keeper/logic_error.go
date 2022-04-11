package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"strconv"
)

// PanicHalt performs an emergency upgrade which immediately halts the chain
// The Team has to come up with a solution and develop a patch to handle
// the update.
// Use this method instead of go's panic() to recover more easily from panics.
// It also leaves the api and rpc end points available.
func (k Keeper) PanicHalt(ctx sdk.Context, message string) {

	// Choose next block for the upgrade
	upgradeBlockHeight := ctx.BlockHeader().Height + 1

	// Create emergency plan
	plan := upgradeTypes.Plan{
		Name:   "emergency_" + strconv.FormatInt(upgradeBlockHeight, 10),
		Height: upgradeBlockHeight,
		Info:   "Emergency Halt; panic occurred; Error:" + message,
	}

	// Directly submit emergency plan
	// Errors can't occur with the current sdk-version
	err := k.upgradeKeeper.ScheduleUpgrade(ctx, plan)
	if err != nil {
		// Can't happen with current sdk
		panic("Emergency Halt failed: " + message)
	}
}
