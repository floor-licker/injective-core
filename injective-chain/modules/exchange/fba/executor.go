package fba

import (
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/InjectiveLabs/metrics/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/keeper/derivative"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/keeper/spot"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

var _ Orderbook = &spot.SpotMarketOrderbook{}
var _ Orderbook = &derivative.MarketOrderbook{}
var _ Orderbook = &spot.SpotLimitOrderbook{}
var _ Orderbook = &derivative.LimitOrderbook{}

type Orderbook interface {
	Peek(ctx sdk.Context) *v2.PriceLevel
	Fill(ctx sdk.Context, quantity math.LegacyDec)
	GetTotalQuantityFilled() math.LegacyDec
	Close() error
}

type ClearingPriceStrategy interface {
	Calculate(lastBuyPrice, lastSellPrice, clearingQuantity math.LegacyDec) math.LegacyDec
}

type BatchAuctionExecutor struct {
	Logger log.Logger
	meter  metrics.Meter
}

func NewBatchAuctionExecutor(logger log.Logger, meter metrics.Meter) *BatchAuctionExecutor {
	e := &BatchAuctionExecutor{
		Logger: logger,
		meter:  meter,
	}

	if e.Logger == nil {
		e.Logger = log.NewNopLogger()
	}

	if e.meter == nil {
		e.meter = metrics.NewNilMeter()
	}

	return e
}

// Match performs the matching algorithm between buy and sell orderbooks.
// Both orderbooks must be non-nil - callers should check this before invoking Match.
func (e *BatchAuctionExecutor) Match(ctx sdk.Context, buyOrderbook, sellOrderbook Orderbook, clearingStrategy ClearingPriceStrategy) *MatchResult {
	defer e.meter.FuncTiming(&ctx, "BatchAuctionExecutor.Execute")()

	lastBuyPrice := math.LegacyZeroDec()
	lastSellPrice := math.LegacyZeroDec()

	for {
		buyOrder := buyOrderbook.Peek(ctx)
		sellOrder := sellOrderbook.Peek(ctx)

		// Base Case: Finished iterating over all the orders
		if buyOrder == nil || sellOrder == nil {
			break
		}

		unitSpread := sellOrder.Price.Sub(buyOrder.Price)
		matchQuantityIncrement := math.LegacyMinDec(buyOrder.Quantity, sellOrder.Quantity)

		// Exit if no more matchable orders
		if unitSpread.IsPositive() || matchQuantityIncrement.IsZero() {
			break
		}

		lastBuyPrice = buyOrder.Price
		lastSellPrice = sellOrder.Price

		buyOrderbook.Fill(ctx, matchQuantityIncrement)
		sellOrderbook.Fill(ctx, matchQuantityIncrement)
	}

	clearingQuantity := sellOrderbook.GetTotalQuantityFilled()
	clearingPrice := clearingStrategy.Calculate(lastBuyPrice, lastSellPrice, clearingQuantity)

	return &MatchResult{
		LastBuyPrice:     lastBuyPrice,
		LastSellPrice:    lastSellPrice,
		ClearingQuantity: clearingQuantity,
		ClearingPrice:    clearingPrice,
	}
}

// EmptyMatchResult returns a MatchResult with zero values, used when orderbooks are nil.
func EmptyMatchResult() *MatchResult {
	return &MatchResult{
		LastBuyPrice:     math.LegacyZeroDec(),
		LastSellPrice:    math.LegacyZeroDec(),
		ClearingQuantity: math.LegacyZeroDec(),
		ClearingPrice:    math.LegacyZeroDec(),
	}
}
