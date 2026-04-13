package base

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

func (k *BaseKeeper) GetMarketAtomicExecutionFeeMultiplier(
	ctx sdk.Context,
	marketId common.Hash,
	marketType types.MarketType,
) math.LegacyDec {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetMarketAtomicExecutionFeeMultiplier")()

	store := k.getStore(ctx)
	takerFeeStore := prefix.NewStore(store, types.AtomicMarketOrderTakerFeeMultiplierKey)

	bz := takerFeeStore.Get(marketId.Bytes())

	if bz == nil {
		return k.GetDefaultAtomicMarketOrderFeeMultiplier(ctx, marketType)
	}

	var multiplier v2.MarketFeeMultiplier
	k.cdc.MustUnmarshal(bz, &multiplier)

	return multiplier.FeeMultiplier
}

// GetDefaultAtomicMarketOrderFeeMultiplier returns the default atomic orders taker fee multiplier for a given market type
func (k *BaseKeeper) GetDefaultAtomicMarketOrderFeeMultiplier(ctx sdk.Context, marketType types.MarketType) math.LegacyDec {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetDefaultAtomicMarketOrderFeeMultiplier")()

	params := k.GetParams(ctx)

	switch marketType {
	case types.MarketType_Spot:
		return params.SpotAtomicMarketOrderFeeMultiplier
	case types.MarketType_Expiry, types.MarketType_Perpetual:
		return params.DerivativeAtomicMarketOrderFeeMultiplier
	case types.MarketType_BinaryOption:
		return params.BinaryOptionsAtomicMarketOrderFeeMultiplier
	default:
		return math.LegacyDec{}
	}
}

func (k *BaseKeeper) GetAllMarketAtomicExecutionFeeMultipliers(ctx sdk.Context) []*v2.MarketFeeMultiplier {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetAllMarketAtomicExecutionFeeMultipliers")()

	store := k.getStore(ctx)
	takerFeeStore := prefix.NewStore(store, types.AtomicMarketOrderTakerFeeMultiplierKey)

	multipliers := make([]*v2.MarketFeeMultiplier, 0)

	iterateSafe(takerFeeStore.Iterator(nil, nil), func(_, value []byte) bool {
		var multiplier v2.MarketFeeMultiplier
		k.cdc.MustUnmarshal(value, &multiplier)
		multipliers = append(multipliers, &multiplier)
		return false
	})

	return multipliers
}

func (k *BaseKeeper) SetAtomicMarketOrderFeeMultipliers(ctx sdk.Context, marketFeeMultipliers []*v2.MarketFeeMultiplier) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetAtomicMarketOrderFeeMultipliers")()

	store := k.getStore(ctx)
	takerFeeStore := prefix.NewStore(store, types.AtomicMarketOrderTakerFeeMultiplierKey)

	for _, multiplier := range marketFeeMultipliers {
		marketID := common.HexToHash(multiplier.MarketId)
		bz := k.cdc.MustMarshal(multiplier)
		takerFeeStore.Set(marketID.Bytes(), bz)
	}
}

func (k *BaseKeeper) AppendOrderExpirations(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
	order *v2.OrderData,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "AppendOrderExpirations")()

	store := k.getStore(ctx)
	expirationStore := prefix.NewStore(store, types.GetOrderExpirationPrefix(expirationBlock, marketID))

	bz := k.cdc.MustMarshal(order)
	expirationStore.Set(common.HexToHash(order.OrderHash).Bytes(), bz)

	expirationMarketsStore := prefix.NewStore(store, types.GetOrderExpirationMarketPrefix(expirationBlock))
	expirationMarketsStore.Set(marketID.Bytes(), []byte{types.TrueByte})
}

func (k *BaseKeeper) DeleteMarketWithOrderExpirations(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteMarketWithOrderExpirations")()

	store := k.getStore(ctx)
	expirationMarketsStore := prefix.NewStore(store, types.GetOrderExpirationMarketPrefix(expirationBlock))
	expirationMarketsStore.Delete(marketID.Bytes())
}

func (k *BaseKeeper) DeleteOrderExpiration(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
	orderHash common.Hash,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteOrderExpiration")()

	store := k.getStore(ctx)
	expirationStore := prefix.NewStore(store, types.GetOrderExpirationPrefix(expirationBlock, marketID))
	expirationStore.Delete(orderHash.Bytes())
}

// DeleteOrderExpirationByKey removes a scheduled order expiration using the raw
// order-hash store key.
func (k *BaseKeeper) DeleteOrderExpirationByKey(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
	orderHashKey []byte,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteOrderExpirationByKey")()

	store := k.getStore(ctx)
	expirationStore := prefix.NewStore(store, types.GetOrderExpirationPrefix(expirationBlock, marketID))
	expirationStore.Delete(orderHashKey)
}

// GetMarketsWithOrderExpirations retrieves all markets with orders expiring at a given block
func (k *BaseKeeper) GetMarketsWithOrderExpirations(
	ctx sdk.Context,
	expirationBlock int64,
) []common.Hash {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetMarketsWithOrderExpirations")()

	store := k.getStore(ctx)
	expirationMarketsStore := prefix.NewStore(store, types.GetOrderExpirationMarketPrefix(expirationBlock))

	markets := make([]common.Hash, 0)

	iterateSafe(expirationMarketsStore.Iterator(nil, nil), func(key, _ []byte) bool {
		marketID := common.BytesToHash(key)
		markets = append(markets, marketID)
		return false
	})

	return markets
}

// GetOrdersByExpiration retrieves all derivative limit orders expiring at a specific block for a market
func (k *BaseKeeper) GetOrdersByExpiration(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
) ([]*v2.OrderData, error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetOrdersByExpiration")()

	orders := make([]*v2.OrderData, 0)

	var err error

	k.IterateOrderExpirationEntries(ctx, marketID, expirationBlock, func(_, value []byte) bool {
		order, unmarshalErr := k.UnmarshalOrderData(value)
		if unmarshalErr != nil {
			err = unmarshalErr
			return true
		}
		orders = append(orders, &order)
		return false
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// UnmarshalOrderData decodes stored order metadata without panicking on malformed bytes.
func (k *BaseKeeper) UnmarshalOrderData(bz []byte) (v2.OrderData, error) {
	var order v2.OrderData
	if err := k.cdc.Unmarshal(bz, &order); err != nil {
		return v2.OrderData{}, err
	}

	return order, nil
}

// IterateOrderExpirationEntries streams raw scheduled order expiration entries
// for a market and block height without materializing them into a slice first.
func (k *BaseKeeper) IterateOrderExpirationEntries(
	ctx sdk.Context,
	marketID common.Hash,
	expirationBlock int64,
	process func(orderHashKey []byte, value []byte) (stop bool),
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "IterateOrderExpirationEntries")()

	store := k.getStore(ctx)
	expirationStore := prefix.NewStore(store, types.GetOrderExpirationPrefix(expirationBlock, marketID))
	iterateSafe(expirationStore.Iterator(nil, nil), process)
}

func (k *BaseKeeper) GetAllMarketIDsWithQuoteDenoms(ctx sdk.Context) []*v2.MarketIDQuoteDenomMakerFee {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetAllMarketIDsWithQuoteDenoms")()

	derivativeMarkets := k.GetAllDerivativeMarkets(ctx)
	spotMarkets := k.GetAllSpotMarkets(ctx)
	binaryOptionsMarkets := k.GetAllBinaryOptionsMarkets(ctx)

	marketIDQuoteDenoms := make([]*v2.MarketIDQuoteDenomMakerFee, 0, len(derivativeMarkets)+len(spotMarkets)+len(binaryOptionsMarkets))

	for _, m := range derivativeMarkets {
		marketIDQuoteDenoms = append(marketIDQuoteDenoms, &v2.MarketIDQuoteDenomMakerFee{
			MarketID:   common.HexToHash(m.MarketId),
			QuoteDenom: m.QuoteDenom,
			MakerFee:   m.MakerFeeRate,
		})
	}

	for _, m := range spotMarkets {
		marketIDQuoteDenoms = append(marketIDQuoteDenoms, &v2.MarketIDQuoteDenomMakerFee{
			MarketID:   m.MarketID(),
			QuoteDenom: m.QuoteDenom,
			MakerFee:   m.MakerFeeRate,
		})
	}

	for _, m := range binaryOptionsMarkets {
		marketIDQuoteDenoms = append(marketIDQuoteDenoms, &v2.MarketIDQuoteDenomMakerFee{
			MarketID:   m.MarketID(),
			QuoteDenom: m.QuoteDenom,
			MakerFee:   m.MakerFeeRate,
		})
	}

	return marketIDQuoteDenoms
}
