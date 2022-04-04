package cli

import (
	"strconv"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdProposalByHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-by-height [pool-id] [height]",
		Short: "Query proposal-by-height",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqHeight, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryProposalByHeightRequest{

				PoolId: reqPoolId,
				Height: reqHeight,
			}

			res, err := queryClient.ProposalByHeight(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
