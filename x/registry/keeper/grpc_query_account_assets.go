package keeper

import (
	"bytes"
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountAssets returns an overview of the balances of the given user regarding the protocol nodes
// This includes the current balance, current staking, delegation, funding and unbondings.
// Supports Pagination
func (k Keeper) AccountAssets(goCtx context.Context, req *types.QueryAccountAssetsRequest) (*types.QueryAccountAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryAccountAssetsResponse{
		Balance:                     0,
		ProtocolStaking:             0,
		ProtocolDelegation:          0,
		ProtocolStakingUnbonding:    0,
		ProtocolDelegationUnbonding: 0,
		ProtocolRewards:             0,
		ProtocolFunding:             0,
	}

	// Fetch account balance
	account, _ := sdk.AccAddressFromBech32(req.Address)
	balance := k.bankKeeper.GetBalance(ctx, account, "tkyve")
	response.Balance = balance.Amount.Uint64()

	// Iterate all Delegator entries
	// Fetches the total delegation and calculates the outstanding rewards
	// TODO find solution to increase performance
	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegatorKeyPrefix))
	delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, nil)

	defer delegatorIterator.Close()

	for ; delegatorIterator.Valid(); delegatorIterator.Next() {
		if bytes.Compare(delegatorIterator.Key()[53:96], []byte(req.Address)) != 0 {
			continue
		}

		var val types.Delegator
		k.cdc.MustUnmarshal(delegatorIterator.Value(), &val)

		f1 := F1Distribution{
			k:                k,
			ctx:              ctx,
			poolId:           val.Id,
			stakerAddress:    val.Staker,
			delegatorAddress: val.Delegator,
		}

		response.ProtocolRewards += f1.getCurrentReward()
		response.ProtocolDelegation += val.DelegationAmount
	}

	// Iterate all Staker entries
	// Fetches the total delegation and calculates the outstanding rewards
	stakerStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	var stakerPrefix []byte
	stakerPrefix = append(stakerPrefix, []byte(req.Address)...)
	stakerPrefix = append(stakerPrefix, []byte("/")...)
	stakerIterator := sdk.KVStorePrefixIterator(stakerStore, stakerPrefix)

	defer stakerIterator.Close()

	for ; stakerIterator.Valid(); stakerIterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(stakerIterator.Value(), &val)

		response.ProtocolStaking += val.Amount
		response.ProtocolStakingUnbonding += val.UnbondingAmount
	}

	// Iterate all Unbonding entries indexed by user to fetch current unbondings
	unbondingStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UnbondingEntriesKeyPrefixByDelegator))
	unbondingIterator := sdk.KVStorePrefixIterator(unbondingStore, []byte(req.Address))

	defer unbondingIterator.Close()

	for ; unbondingIterator.Valid(); unbondingIterator.Next() {
		var val types.UnbondingEntries
		k.cdc.MustUnmarshal(unbondingIterator.Value(), &val)

		response.ProtocolDelegationUnbonding += val.Amount
	}

	// Correct value as the unbonding Iterator contains both types of unbonding
	response.ProtocolDelegationUnbonding -= response.ProtocolStakingUnbonding

	// Iterate all funding entries
	funderStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))
	funderIterator := sdk.KVStorePrefixIterator(funderStore, []byte(req.Address))

	defer funderIterator.Close()

	for ; funderIterator.Valid(); funderIterator.Next() {
		var val types.Funder
		k.cdc.MustUnmarshal(funderIterator.Value(), &val)

		response.ProtocolFunding += val.Amount
	}

	return &response, nil
}
