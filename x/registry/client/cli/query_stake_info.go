package cli

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdStakeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-info [pool-id] [staker] [current-stake] [minimum-stake]",
		Short: "Query stake_info",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			reqStaker := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryStakeInfoRequest{
				PoolId: reqPoolId,
				Staker: reqStaker,
			}

			res, err := queryClient.StakeInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
