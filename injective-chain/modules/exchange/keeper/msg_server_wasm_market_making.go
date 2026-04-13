package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

type WasmMsgServer struct {
	*Keeper
}

// NewWasmMsgServerImpl returns an implementation of the exchange MsgServer interface for the provided Keeper for exchange wasm functions.
func NewWasmMsgServerImpl(keeper *Keeper) WasmMsgServer {
	return WasmMsgServer{
		Keeper: keeper,
	}
}

func (k WasmMsgServer) PrivilegedExecuteContract(
	c context.Context,
	msg *v2.MsgPrivilegedExecuteContract,
) (*v2.MsgPrivilegedExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "PrivilegedExecuteContract")()

	return k.PrivilegedExecuteContractWithVersion(ctx, msg, types.ExchangeTypeVersionV2)
}
