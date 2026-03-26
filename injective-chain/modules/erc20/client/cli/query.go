package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/InjectiveLabs/injective-core/cli"
	cliflags "github.com/InjectiveLabs/injective-core/cli/flags"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/erc20/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := cli.ModuleRootCommand(types.ModuleName, true)

	cmd.AddCommand(
		GetParams(),
		GetTokenPairs(),
		GetTokenPairByDenom(),
		GetTokenPairByERC20(),
	)

	return cmd
}

func GetParams() *cobra.Command {
	return cli.QueryCmd("params",
		"Gets module params",
		types.NewQueryClient,
		&types.QueryParamsRequest{}, nil, nil,
	)
}

func GetTokenPairs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-pairs",
		Short: "Returns all created token pairs in the module",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllTokenPairs(cmd.Context(), &types.QueryAllTokenPairsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cliflags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "token-pairs")

	return cmd
}

func GetTokenPairByDenom() *cobra.Command {
	return cli.QueryCmd("token-pair-by-denom <denom>",
		"Returns the token pair associated with denom",
		types.NewQueryClient,
		&types.QueryTokenPairByDenomRequest{}, nil, nil,
	)
}

func GetTokenPairByERC20() *cobra.Command {
	return cli.QueryCmd("token-pair-by-erc20 <erc20 address>",
		"Returns the token pair associated with erc20 address",
		types.NewQueryClient,
		&types.QueryTokenPairByERC20AddressRequest{}, nil, nil,
	)
}
