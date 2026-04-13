package wasmx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/keeper"
)

type BlockHandler struct {
	k *keeper.Keeper
}

func NewBlockHandler(k keeper.Keeper) *BlockHandler {
	return &BlockHandler{
		k: &k,
	}
}
func (h *BlockHandler) BeginBlocker(ctx sdk.Context) error {
	defer h.k.Meter(ctx).FuncTiming(&ctx, "BeginBlocker")()

	return h.k.ExecuteContracts(ctx)
}
