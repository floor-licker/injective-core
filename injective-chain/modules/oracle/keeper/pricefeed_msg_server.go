package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

type PricefeedMsgServer struct {
	*Keeper
}

// NewPricefeedMsgServerImpl returns an implementation of the price feed provider MsgServer interface for the provided Keeper for price feed provider oracle functions.
func NewPricefeedMsgServerImpl(keeper Keeper) PricefeedMsgServer {
	return PricefeedMsgServer{
		Keeper: &keeper,
	}
}

func (k PricefeedMsgServer) RelayPriceFeedPrice(c context.Context, msg *types.MsgRelayPriceFeedPrice) (*types.MsgRelayPriceFeedPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RelayPriceFeedPrice")()
	// prepare context
	if err := k.ProcessPriceFeedPrice(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgRelayPriceFeedPriceResponse{}, nil
}
