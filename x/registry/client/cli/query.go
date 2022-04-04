package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KYVENetwork/chain/x/registry/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group registry queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListPool())
	cmd.AddCommand(CmdShowPool())
	cmd.AddCommand(CmdCanPropose())
	cmd.AddCommand(CmdCanVote())
	cmd.AddCommand(CmdFundingAmount())
	cmd.AddCommand(CmdStakingAmount())

	cmd.AddCommand(CmdFundersList())

	cmd.AddCommand(CmdStakersList())

	cmd.AddCommand(CmdAccountFundedList())

	cmd.AddCommand(CmdAccountStakedList())

	cmd.AddCommand(CmdShowDelegator())

	cmd.AddCommand(CmdListProposal())
	cmd.AddCommand(CmdShowProposal())
	cmd.AddCommand(CmdProposalByHeight())

	cmd.AddCommand(CmdAccountStakersDelegationList())

	cmd.AddCommand(CmdStakersByPoolAndDelegator())

	// this line is used by starport scaffolding # 1

	return cmd
}
