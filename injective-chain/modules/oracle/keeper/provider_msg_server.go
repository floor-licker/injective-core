package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

type ProviderMsgServer struct {
	ProviderKeeper
}

// NewProviderMsgServerImpl returns an implementation of the provider MsgServer interface for the provided Keeper for provider oracle functions.
func NewProviderMsgServerImpl(keeper Keeper) ProviderMsgServer {
	return ProviderMsgServer{
		ProviderKeeper: &keeper,
	}
}

func (k ProviderMsgServer) RelayProviderPrices(c context.Context, msg *types.MsgRelayProviderPrices) (*types.MsgRelayProviderPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RelayProviderPrices")()

	relayer, _ := sdk.AccAddressFromBech32(msg.Sender)
	if !k.IsProviderRelayer(ctx, msg.Provider, relayer) {
		return nil, errors.Wrapf(types.ErrRelayerNotAuthorized, "relayer %s not an authorized provider for %s", relayer.String(), msg.Provider)
	}

	k.ProcessProviderPrices(ctx, msg)
	return &types.MsgRelayProviderPricesResponse{}, nil
}
