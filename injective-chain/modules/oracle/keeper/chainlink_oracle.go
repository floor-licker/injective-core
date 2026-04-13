package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

// GetChainlinkPriceState reads the stored price state.
func (k *Keeper) GetChainlinkPriceState(ctx sdk.Context, symbol string) *types.ChainlinkPriceState {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetChainlinkPriceState")()

	var priceState types.ChainlinkPriceState
	bz := k.getStore(ctx).Get(types.GetChainlinkPriceStoreKey(symbol))
	if bz == nil {
		return nil
	}

	k.cdc.MustUnmarshal(bz, &priceState)
	return &priceState
}

// GetAllChainlinkPriceStates reads all stored chainlink price states.
func (k *Keeper) GetAllChainlinkPriceStates(ctx sdk.Context) []*types.ChainlinkPriceState {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetAllChainlinkPriceStates")()

	priceStates := make([]*types.ChainlinkPriceState, 0)
	store := ctx.KVStore(k.storeKey)
	chainlinkPriceStore := prefix.NewStore(store, types.ChainlinkPriceKey)

	iterator := chainlinkPriceStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var chainlinkPriceState types.ChainlinkPriceState
		k.cdc.MustUnmarshal(iterator.Value(), &chainlinkPriceState)
		priceStates = append(priceStates, &chainlinkPriceState)
	}

	return priceStates
}
