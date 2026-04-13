package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
)

func (k *Keeper) GetHighestBid(ctx sdk.Context) *types.Bid {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetHighestBid")()

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BidsKey)

	if bz == nil {
		return &types.Bid{
			Bidder: "",
			Amount: sdk.NewCoin("inj", math.ZeroInt()),
		}
	}

	var bid types.Bid
	k.cdc.MustUnmarshal(bz, &bid)
	return &bid
}

func (k *Keeper) SetBid(ctx sdk.Context, sender string, amount sdk.Coin) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetBid")()

	store := ctx.KVStore(k.storeKey)
	bid := &types.Bid{
		Bidder: sender,
		Amount: amount,
	}
	bz := k.cdc.MustMarshal(bid)
	store.Set(types.BidsKey, bz)
}

func (k *Keeper) DeleteBid(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteBid")()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.BidsKey)
}
