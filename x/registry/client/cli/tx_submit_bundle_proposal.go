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
		Use:   "submit-bundle-proposal [id] [bundle-id] [byte-size] [from-height] [to-height] [from-key] [to-key] [to-value]",
		Short: "Broadcast message submit-bundle-proposal",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argStorageId := args[1]
			argByteSize, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argFromHeight, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argToHeight, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argFromKey := args[5]
			argToKey := args[6]
			argToValue := args[7]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitBundleProposal(
				clientCtx.GetFromAddress().String(),
				argId,
				argStorageId,
				argByteSize,
				argFromHeight,
				argToHeight,
				argFromKey,
				argToKey,
				argToValue,
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
