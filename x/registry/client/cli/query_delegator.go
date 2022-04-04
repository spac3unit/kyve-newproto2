package cli

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdShowDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-delegator [id]",
		Short: "shows a Delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorRequest{
				PoolId: argId,
			}

			res, err := queryClient.Delegator(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
