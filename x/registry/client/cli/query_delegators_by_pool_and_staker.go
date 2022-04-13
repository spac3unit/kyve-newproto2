package cli

import (
	"github.com/spf13/cast"
	"strconv"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDelegatorsByPoolAndStaker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegators-by-pool-and-staker [pool-id] [staker]",
		Short: "Query account_stakers_delegation_list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqStaker := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegatorsByPoolAndStakerRequest{

				PoolId: reqPoolId,
				Staker: reqStaker,
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params.Pagination = pageReq

			res, err := queryClient.DelegatorsByPoolAndStaker(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
