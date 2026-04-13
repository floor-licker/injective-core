package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	osmosistypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/txfees/osmosis/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/txfees/types"
)

var _ types.QueryServer = queryServer{}

// queryServer defines a wrapper around the x/txfees keeper providing gRPC method
// handlers.
type queryServer struct {
	k *Keeper
}

func NewQueryServer(k *Keeper) types.QueryServer {
	return queryServer{
		k: k,
	}
}

func (q queryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.k.Meter(ctx).FuncTiming(&ctx, "Params")()

	params := q.k.GetParams(ctx)

	res := &types.QueryParamsResponse{
		Params: params,
	}

	return res, nil
}

// since we only store current baseFee, this query can only return current BaseFee values, even if requested for historic blocks
func (q queryServer) GetEipBaseFee(c context.Context, _ *types.QueryEipBaseFeeRequest) (*types.QueryEipBaseFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.k.Meter(ctx).FuncTiming(&ctx, "GetEipBaseFee")()

	if ctx.BlockHeight() < q.k.CurFeeState.GetCurrentBlockHeight()-2 { // we do not support historical queries since we only have current BaseFee in memory
		return nil, types.ErrUnsupportedQueryParams
	}

	baseFee := q.k.CurFeeState.GetCurBaseFee()
	return &types.QueryEipBaseFeeResponse{BaseFee: &types.EipBaseFee{BaseFee: baseFee}}, nil
}

var _ osmosistypes.QueryServer = osmosisQueryServer{}

type osmosisQueryServer struct {
	k *Keeper
}

func NewOsmosisQueryServer(k *Keeper) osmosistypes.QueryServer {
	return osmosisQueryServer{
		k: k,
	}
}

func (q osmosisQueryServer) GetEipBaseFee(
	c context.Context, _ *osmosistypes.QueryEipBaseFeeRequest,
) (*osmosistypes.QueryEipBaseFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.k.Meter(ctx).FuncTiming(&ctx, "osmosisQueryServer.GetEipBaseFee")()

	if ctx.BlockHeight() < q.k.CurFeeState.GetCurrentBlockHeight()-2 { // we do not support historical queries since we only have current BaseFee in memory
		return nil, types.ErrUnsupportedQueryParams
	}

	response := q.k.CurFeeState.GetCurBaseFee()
	return &osmosistypes.QueryEipBaseFeeResponse{BaseFee: response}, nil
}
