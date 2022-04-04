package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
)

var _ types.QueryServer = Keeper{}
