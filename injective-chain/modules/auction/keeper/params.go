package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
)

// AuctionPeriodDuration auction period param
func (k *Keeper) AuctionPeriodDuration(ctx sdk.Context) int64 {
	defer k.Meter(ctx).FuncTiming(&ctx, "AuctionPeriodDuration")()

	return k.GetParams(ctx).AuctionPeriod
}

// MinNextBidIncrementRate returns min percentage increment param
func (k *Keeper) MinNextBidIncrementRate(ctx sdk.Context) string {
	defer k.Meter(ctx).FuncTiming(&ctx, "MinNextBidIncrementRate")()

	return k.GetParams(ctx).MinNextBidIncrementRate.String()
}

// GetParams returns the total set of auction parameters.
func (k *Keeper) GetParams(ctx sdk.Context) types.Params {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetParams")()

	store := k.GetStore(ctx)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.Params{}
	}

	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

// SetParams set the params
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetParams")()

	store := k.GetStore(ctx)
	store.Set(types.ParamsKey, k.cdc.MustMarshal(&params))
}
