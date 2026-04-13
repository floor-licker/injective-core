package govcli

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// flagNoIndent matches autocli's FlagNoIndent for parity (do not indent JSON output).
const flagNoIndent = "no-indent"

// AutoCLI-style proposal status values (see x/gov/autocli.go Example).
const (
	proposalStatusUnspecified   = "unspecified"
	proposalStatusDepositPeriod = "deposit-period"
	proposalStatusVotingPeriod  = "voting-period"
	proposalStatusPassed        = "passed"
	proposalStatusRejected      = "rejected"
	proposalStatusFailed        = "failed"
)

// addGovQueryFlags adds query connection flags and --no-indent to cmd (parity with autocli-generated gov queries).
// Keyring flags are not added; they are only used for tx signing and key management, not for read-only queries.
func addGovQueryFlags(cmd *cobra.Command) {
	cmd.Flags().String(flags.FlagChainID, "", "network chain ID")
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().BoolP(flagNoIndent, "", false, "Do not indent JSON output")
}

// proposalStatusFromString accepts AutoCLI-style values (unspecified, deposit-period, ...)
// and optionally PROPOSAL_STATUS_* enum constants.
func proposalStatusFromString(s string) (govv1.ProposalStatus, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return govv1.ProposalStatus_PROPOSAL_STATUS_UNSPECIFIED, nil
	}
	switch strings.ToLower(s) {
	case proposalStatusUnspecified:
		return govv1.ProposalStatus_PROPOSAL_STATUS_UNSPECIFIED, nil
	case proposalStatusDepositPeriod:
		return govv1.ProposalStatus_PROPOSAL_STATUS_DEPOSIT_PERIOD, nil
	case proposalStatusVotingPeriod:
		return govv1.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD, nil
	case proposalStatusPassed:
		return govv1.ProposalStatus_PROPOSAL_STATUS_PASSED, nil
	case proposalStatusRejected:
		return govv1.ProposalStatus_PROPOSAL_STATUS_REJECTED, nil
	case proposalStatusFailed:
		return govv1.ProposalStatus_PROPOSAL_STATUS_FAILED, nil
	}
	return govv1.ProposalStatusFromString(s)
}

// GetQueryCmd returns the custom gov query command with only proposal and proposals
// subcommands. Output uses the app codec (clientCtx.PrintProto) so nested Anys with
// LegacyDec fields (e.g. txfees Params) marshal without panic. Autocli adds the
// remaining gov query subcommands (vote, votes, deposit, deposits, tally, params,
// constitution) via EnhanceCustomCommand and encodes those with aminojson (safe).
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        govtypes.ModuleName,
		Short:                      "Querying commands for the gov module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetProposalCmd(),
		GetProposalsCmd(),
	)

	return cmd
}

// GetProposalCmd returns the query proposal by id command.
func GetProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposal [proposal-id]",
		Aliases: []string{"proposer"},
		Short:   "Query details of a single proposal",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := govv1.NewQueryClient(clientCtx)
			res, err := queryClient.Proposal(cmd.Context(), &govv1.QueryProposalRequest{
				ProposalId: proposalID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	addGovQueryFlags(cmd)
	return cmd
}

// GetProposalsCmd returns the query proposals command with filters and pagination.
// Flag names match cosmos.gov.v1.QueryProposalsRequest field names (proposal_status as proposal-status for CLI), consistent with autocli.
func GetProposalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			depositor, _ := cmd.Flags().GetString("depositor")
			voter, _ := cmd.Flags().GetString("voter")
			statusStr, _ := cmd.Flags().GetString("proposal-status")

			req := &govv1.QueryProposalsRequest{
				Depositor: depositor,
				Voter:     voter,
			}

			if statusStr != "" {
				status, err := proposalStatusFromString(statusStr)
				if err != nil {
					return err
				}
				req.ProposalStatus = status
			}

			fs, err := client.FlagSetWithPageKeyDecoded(cmd.Flags())
			if err != nil {
				return err
			}
			pageReq, err := readPageRequest(fs)
			if err != nil {
				return err
			}
			req.Pagination = pageReq

			queryClient := govv1.NewQueryClient(clientCtx)
			res, err := queryClient.Proposals(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("depositor", "", "filter by depositor address")
	cmd.Flags().String("voter", "", "filter by voter address")
	cmd.Flags().String("proposal-status", "", "filter by proposal status (unspecified|deposit-period|voting-period|passed|rejected|failed or PROPOSAL_STATUS_* enum)")
	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	addAutoCLIPaginationFlags(cmd, "proposals")
	addGovQueryFlags(cmd)
	return cmd
}

// addAutoCLIPaginationFlags adds AutoCLI-style page-* flags (page-limit, page-offset, etc.)
// for parity with autocli-generated commands. page-key is already added by AddPaginationFlagsToCmd as FlagPageKey.
func addAutoCLIPaginationFlags(cmd *cobra.Command, queryName string) {
	cmd.Flags().Uint64("page-offset", 0, "pagination offset of "+queryName+" to query for")
	cmd.Flags().Uint64("page-limit", 100, "pagination limit of "+queryName+" to query for")
	cmd.Flags().Bool("page-count-total", false, "count total number of records in "+queryName+" to query for")
	cmd.Flags().Bool("page-reverse", false, "results are sorted in descending order")
}

// readPageRequest reads pagination from either AutoCLI-style (page-*) or legacy (limit, offset, ...) flags.
func readPageRequest(fs *pflag.FlagSet) (*query.PageRequest, error) {
	var limit, offset uint64
	var countTotal, reverse bool
	var pageKey string
	if fs.Changed("page-limit") {
		limit, _ = fs.GetUint64("page-limit")
	} else {
		limit, _ = fs.GetUint64(flags.FlagLimit)
	}
	if fs.Changed("page-offset") {
		offset, _ = fs.GetUint64("page-offset")
	} else {
		offset, _ = fs.GetUint64(flags.FlagOffset)
	}
	if fs.Changed("page-count-total") {
		countTotal, _ = fs.GetBool("page-count-total")
	} else {
		countTotal, _ = fs.GetBool(flags.FlagCountTotal)
	}
	if fs.Changed("page-reverse") {
		reverse, _ = fs.GetBool("page-reverse")
	} else {
		reverse, _ = fs.GetBool(flags.FlagReverse)
	}
	pageKey, _ = fs.GetString(flags.FlagPageKey)
	page, _ := fs.GetUint64(flags.FlagPage)
	if page > 1 && offset > 0 {
		return nil, errors.New("page and offset cannot be used together")
	}
	if page > 1 {
		offset = (page - 1) * limit
	}
	return &query.PageRequest{
		Key:        []byte(pageKey),
		Offset:     offset,
		Limit:      limit,
		CountTotal: countTotal,
		Reverse:    reverse,
	}, nil
}
