package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

type StorkMsgServer struct {
	*Keeper
}

// NewStorkMsgServerImpl returns an implementation of the stork provider MsgServer interface for the provided Keeper for coinbase provider oracle functions.
func NewStorkMsgServerImpl(keeper Keeper) StorkMsgServer {
	return StorkMsgServer{
		Keeper: &keeper,
	}
}

func (k StorkMsgServer) RelayStorkMessage(c context.Context, msg *types.MsgRelayStorkPrices) (*types.MsgRelayStorkPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RelayStorkMessage")()

	k.ProcessStorkAssetPairsData(ctx, msg.AssetPairs)

	return &types.MsgRelayStorkPricesResponse{}, nil
}
