package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

type PythMsgServer struct {
	*Keeper
}

// NewPythMsgServerImpl returns an implementation of the Pyth MsgServer interface for the provided Keeper for Pyth oracle functions.
func NewPythMsgServerImpl(keeper Keeper) PythMsgServer {
	return PythMsgServer{
		Keeper: &keeper,
	}
}

func (k PythMsgServer) RelayPythPrices(c context.Context, msg *types.MsgRelayPythPrices) (*types.MsgRelayPythPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RelayPythPrices")()

	if len(msg.PriceAttestations) == 0 {
		return &types.MsgRelayPythPricesResponse{}, nil
	}

	pythContract, err := sdk.AccAddressFromBech32(k.GetParams(ctx).PythContract)
	if err != nil {
		return nil, types.ErrPythContractNotFound
	}

	if !pythContract.Equals(sdk.MustAccAddressFromBech32(msg.Sender)) {
		return nil, types.ErrUnauthorizedPythPriceRelay
	}

	k.ProcessPythPriceAttestations(ctx, msg.PriceAttestations)
	return &types.MsgRelayPythPricesResponse{}, nil
}
