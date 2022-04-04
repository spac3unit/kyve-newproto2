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

func CmdUpdateMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-metadata [id] [commission] [moniker] [website] [logo]",
		Short: "Broadcast message update-metadata",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argCommission := args[1]
			moniker := args[2]
			website := args[3]
			logo := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateMetadata(
				clientCtx.GetFromAddress().String(),
				argId,
				argCommission,
				moniker,
				website,
				logo,
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
