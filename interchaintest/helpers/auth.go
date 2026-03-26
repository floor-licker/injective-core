package helpers

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// QueryProposalRPC queries a proposal via gRPC
func QueryModuleAccount(ctx context.Context, chain *cosmos.CosmosChain, module string) (sdk.AccountI, error) {
	conn, err := grpc.NewClient(chain.GetHostGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryClient := auth.NewQueryClient(conn)

	resp, err := QueryRPC(ctx, queryClient.ModuleAccountByName, &auth.QueryModuleAccountByNameRequest{
		Name: module,
	})
	if err != nil {
		return nil, err
	}

	cdc := chain.GetCodec()
	auth.RegisterInterfaces(cdc.InterfaceRegistry())

	var acc sdk.AccountI
	if err := chain.GetCodec().UnpackAny(resp.Account, &acc); err != nil {
		return nil, err
	}

	return acc, nil
}
