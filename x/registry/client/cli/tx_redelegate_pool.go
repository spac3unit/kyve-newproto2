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

func CmdRedelegatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate-pool [from_pool_id] [from_staker] [to_pool_id] [to_staker] [amount]",
		Short: "Broadcast message redelegate-pool",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			fromPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			fromStaker := args[1]

			toPoolId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			toStaker := args[3]

			amount, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRedelegatePool{
				Creator:    clientCtx.GetFromAddress().String(),
				FromPoolId: fromPoolId,
				FromStaker: fromStaker,
				ToPoolId:   toPoolId,
				ToStaker:   toStaker,
				Amount:     amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
