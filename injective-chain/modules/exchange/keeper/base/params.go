package base

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

// GetParams returns the total set of exchange parameters.
func (k *BaseKeeper) GetParams(ctx sdk.Context) v2.Params {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetParams")()

	return k.GetCachedParams(ctx)
}

func (k *BaseKeeper) getParamsFromStore(ctx sdk.Context) v2.Params {
	store := k.getStore(ctx)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return v2.Params{}
	}

	var params v2.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

// SetParams set the params and updates the object cache
func (k *BaseKeeper) SetParams(ctx sdk.Context, params v2.Params) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetParams")()

	store := k.getStore(ctx)
	store.Set(types.ParamsKey, k.cdc.MustMarshal(&params))

	// Update the object cache with the new params
	k.SetCachedParams(ctx, params)
}

func (k *BaseKeeper) IsPostOnlyMode(ctx sdk.Context) bool {
	return k.GetCachedParams(ctx).PostOnlyModeHeightThreshold > ctx.BlockHeight()
}

func (k *BaseKeeper) GetMinimalProtocolFeeRate(ctx sdk.Context, market v2.MarketI) math.LegacyDec {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetMinimalProtocolFeeRate")()

	if market.GetDisabledMinimalProtocolFee() {
		return math.LegacyZeroDec()
	}

	return k.GetCachedParams(ctx).MinimalProtocolFeeRate
}
