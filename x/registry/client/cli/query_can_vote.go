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

func CmdCanVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "can-vote [id] [voter]",
		Short: "Query canVote",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqBundleId := args[1]
			reqVoter := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCanVoteRequest{

				PoolId:   reqId,
				BundleId: reqBundleId,
				Voter:    reqVoter,
			}

			res, err := queryClient.CanVote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
