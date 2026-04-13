package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/peggy/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

func (k *Keeper) GetMintAmounts(ctx sdk.Context) []*types.MintAmount {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetMintAmounts")()

	var mintAmounts []*types.MintAmount
	mintAmountsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.MintAmountERC20Key)
	chaintypes.IterateSafe(mintAmountsStore.Iterator(nil, nil), func(k, v []byte) (stop bool) {
		var amount sdkmath.Int
		if err := amount.Unmarshal(v); err != nil {
			panic(err)
		}

		mintAmounts = append(mintAmounts, &types.MintAmount{
			Token:  gethcommon.BytesToAddress(k).Hex(),
			Amount: amount,
		})

		return false
	})

	return mintAmounts
}

func (k *Keeper) GetMintAmountERC20(ctx sdk.Context, tokenAddress gethcommon.Address) sdkmath.Int {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetMintAmountERC20")()

	store := k.getStore(ctx)
	bz := store.Get(types.GetMintAmountERC20Key(tokenAddress.Bytes()))
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	}

	var amount sdkmath.Int
	if err := amount.Unmarshal(bz); err != nil {
		panic(err)
	}

	return amount
}

func (k *Keeper) SetMintAmountERC20(ctx sdk.Context, tokenAddress gethcommon.Address, amount sdkmath.Int) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetMintAmountERC20")()

	store := k.getStore(ctx)
	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(types.GetMintAmountERC20Key(tokenAddress.Bytes()), bz)
}

func (k *Keeper) DeleteMintAmountERC20(ctx sdk.Context, tokenAddress gethcommon.Address) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteMintAmountERC20")()

	store := k.getStore(ctx)
	store.Delete(types.GetMintAmountERC20Key(tokenAddress.Bytes()))
}
