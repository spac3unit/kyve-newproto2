package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountAssets(goCtx context.Context, req *types.QueryAccountAssetsRequest) (*types.QueryAccountAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	account, _ := sdk.AccAddressFromBech32(req.Address)
	balance := k.bankKeeper.GetBalance(ctx, account, "tkyve")

	// Fetch all Delegation entries for Delegator with requested address
	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DelegatorKeyPrefix))
	var delegatorKeyPrefix []byte
	delegatorKeyPrefix = append(delegatorKeyPrefix, []byte(req.Address)...)
	delegatorKeyPrefix = append(delegatorKeyPrefix, []byte("/")...)
	delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, delegatorKeyPrefix)

	defer delegatorIterator.Close()

	var protocolDelegation uint64 = 0
	var protocolUnbonding uint64 = 0
	var protocolRewards uint64 = 0

	for ; delegatorIterator.Valid(); delegatorIterator.Next() {
		var val types.Delegator
		k.cdc.MustUnmarshal(delegatorIterator.Value(), &val)

		f1 := F1Distribution{
			k:                k,
			ctx:              ctx,
			poolId:           val.Id,
			stakerAddress:    val.Staker,
			delegatorAddress: val.Delegator,
		}

		protocolRewards += f1.getCurrentReward()
		protocolDelegation += val.DelegationAmount
	}

	// Fetch all Staker entries for requested address
	stakerStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	var stakerPrefix []byte
	stakerPrefix = append(stakerPrefix, []byte(req.Address)...)
	stakerPrefix = append(stakerPrefix, []byte("/")...)
	stakerIterator := sdk.KVStorePrefixIterator(stakerStore, stakerPrefix)

	defer stakerIterator.Close()

	var protocolStaking uint64 = 0

	for ; stakerIterator.Valid(); stakerIterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(stakerIterator.Value(), &val)

		protocolStaking += val.Amount
		protocolUnbonding += val.UnbondingAmount
	}

	return &types.QueryAccountAssetsResponse{
		Balance:             balance.Amount.Uint64(),
		ProtocolStaking:     protocolStaking,
		ProtocolDelegation:  protocolDelegation,
		ProtocolUnbonding:   protocolUnbonding,
		ProtocolRewards:     protocolRewards,
		ValidatorStaking:    0,
		ValidatorDelegation: 0,
		ValidatorUnbonding:  0,
	}, nil
}
