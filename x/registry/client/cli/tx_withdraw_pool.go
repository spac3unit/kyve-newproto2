package cli

import (
	"strconv"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdWithdrawPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-pool [pool-id] [staker-address]",
		Short: "Broadcast message withdraw-pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argStakerAddress := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawPool(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argStakerAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
