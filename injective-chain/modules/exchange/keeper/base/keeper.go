package base

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/InjectiveLabs/metrics/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

//nolint:revive // ok
type BaseKeeper struct {
	storeKey       storetypes.StoreKey
	tStoreKey      storetypes.StoreKey
	objectStoreKey storetypes.StoreKey
	cdc            codec.BinaryCodec
	meter          metrics.Meter
}

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	tStoreKey storetypes.StoreKey,
	objectStoreKey storetypes.StoreKey,
) *BaseKeeper {
	return &BaseKeeper{
		storeKey:       storeKey,
		tStoreKey:      tStoreKey,
		objectStoreKey: objectStoreKey,
		cdc:            cdc,
	}
}

type whiteKnightLiquidatorsSet map[string]struct{}

func (k *BaseKeeper) Meter(ctx context.Context) metrics.Meter {
	if k.meter == nil {
		ctxMeter := sdk.UnwrapSDKContext(ctx).Meter()
		if ctxMeter == nil { // in tests with empty sdk.Context{}
			ctxMeter = metrics.NewNilMeter()
		}
		k.meter = ctxMeter.SubMeter(types.ModuleName, metrics.Tag("svc", types.ModuleName))
	}

	return k.meter
}

func (k *BaseKeeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k *BaseKeeper) GetCodec() codec.BinaryCodec {
	return k.cdc
}

func (k *BaseKeeper) getStore(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k *BaseKeeper) getTransientStore(ctx sdk.Context) storetypes.KVStore {
	return ctx.TransientStore(k.tStoreKey)
}

func (k *BaseKeeper) getObjectStore(ctx sdk.Context) storetypes.ObjKVStore {
	return ctx.ObjectStore(k.objectStoreKey)
}

// SetPostOnlyModeCancellationFlag sets a flag in the store to indicate that post-only mode
// should be cancelled in the next BeginBlock
func (k *BaseKeeper) SetPostOnlyModeCancellationFlag(ctx sdk.Context) {
	store := k.getStore(ctx)
	store.Set(types.PostOnlyModeCancellationKey, []byte{1})
}

// HasPostOnlyModeCancellationFlag checks if the post-only mode cancellation flag is set
func (k *BaseKeeper) HasPostOnlyModeCancellationFlag(ctx sdk.Context) bool {
	store := k.getStore(ctx)
	return store.Has(types.PostOnlyModeCancellationKey)
}

// DeletePostOnlyModeCancellationFlag removes the post-only mode cancellation flag from the store
func (k *BaseKeeper) DeletePostOnlyModeCancellationFlag(ctx sdk.Context) {
	store := k.getStore(ctx)
	store.Delete(types.PostOnlyModeCancellationKey)
}

// ShouldEmitLegacyVersionEvents returns the cached EmitLegacyVersionEvents param value.
// Reads from object store cache, falling back to main store if not cached.
func (k *BaseKeeper) ShouldEmitLegacyVersionEvents(ctx sdk.Context) bool {
	return k.GetCachedParams(ctx).EmitLegacyVersionEvents
}

// SetCachedParams caches the params in the object store for the current block.
// Uses object store to avoid repeated marshal/unmarshal while keeping CacheContext isolation.
// Cached params are treated as immutable for the lifetime of the context.
func (k *BaseKeeper) SetCachedParams(ctx sdk.Context, params v2.Params) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetCachedParams")()

	objStore := k.getObjectStore(ctx)
	paramsCopy := params
	objStore.Set(types.ObjectCachedParamsKey, &paramsCopy)
	objStore.Set(types.ObjectCachedWhiteKnightLiquidatorsKey, buildWhiteKnightLiquidatorsSet(params.WhiteKnightLiquidators))
}

// GetCachedParams returns cached params from object store if available, otherwise fetches from main store.
// This avoids repeated KVStore reads and unmarshalling within the same block.
// Returned params must be treated as immutable by callers.
func (k *BaseKeeper) GetCachedParams(ctx sdk.Context) v2.Params {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetCachedParams")()

	objStore := k.getObjectStore(ctx)
	cached := objStore.Get(types.ObjectCachedParamsKey)
	if cached != nil {
		if params, ok := cached.(*v2.Params); ok {
			return *params
		}
	}

	// Fallback to store read if cache miss (e.g., in genesis, tests, or CheckTx)
	// Populate cache for subsequent calls within the same context
	params := k.getParamsFromStore(ctx)
	k.SetCachedParams(ctx, params)
	return params
}

func (k *BaseKeeper) IsWhiteKnightLiquidator(ctx sdk.Context, address string) bool {
	objStore := k.getObjectStore(ctx)
	cached := objStore.Get(types.ObjectCachedWhiteKnightLiquidatorsKey)
	if cachedSet, ok := cached.(whiteKnightLiquidatorsSet); ok {
		_, found := cachedSet[address]
		return found
	}

	// Fallback path when only params cache exists without whitelist set cache.
	params := k.GetCachedParams(ctx)
	cachedSet := buildWhiteKnightLiquidatorsSet(params.WhiteKnightLiquidators)
	objStore.Set(types.ObjectCachedWhiteKnightLiquidatorsKey, cachedSet)

	_, found := cachedSet[address]
	return found
}

func buildWhiteKnightLiquidatorsSet(whiteKnightLiquidators []string) whiteKnightLiquidatorsSet {
	cachedSet := make(whiteKnightLiquidatorsSet, len(whiteKnightLiquidators))
	for _, whiteKnightLiquidator := range whiteKnightLiquidators {
		whiteKnightLiquidatorAddr, err := sdk.AccAddressFromBech32(whiteKnightLiquidator)
		if err != nil {
			continue
		}
		cachedSet[whiteKnightLiquidatorAddr.String()] = struct{}{}
	}
	return cachedSet
}
