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

func CmdSubmitBundleProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-bundle-proposal [id] [bundle-id] [byte-size] [bundle-size]",
		Short: "Broadcast message submit-bundle-proposal",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argBundleId := args[1]
			argByteSize, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argFromHeight, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argBundleSize, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitBundleProposal(
				clientCtx.GetFromAddress().String(),
				argId,
				argBundleId,
				argByteSize,
				argFromHeight,
				argBundleSize,
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
