package fba

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/keeper/derivative"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/keeper/spot"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

// BatchAuction coordinates FBA (Frequent Batch Auction) execution across all markets.
// It handles both market orders (which match against resting limit orders) and
// limit orders (which match against each other in a batch auction).
type BatchAuction struct {
	spotKeeper       spot.SpotKeeper
	derivativeKeeper derivative.DerivativeKeeper
}

// NewBatchAuction creates a new batch auction coordinator.
func NewBatchAuction(spotKeeper spot.SpotKeeper, derivativeKeeper derivative.DerivativeKeeper) *BatchAuction {
	return &BatchAuction{
		spotKeeper:       spotKeeper,
		derivativeKeeper: derivativeKeeper,
	}
}

// ExecuteMarketOrders executes all market orders (spot and derivative) in parallel.
// Market orders match against resting limit orders in the orderbook.
// Returns the execution results for spot and derivative markets.
func (ba *BatchAuction) ExecuteMarketOrders(
	ctx sdk.Context,
	spotMarketOrderIndicators []*v2.MarketOrderIndicator,
	derivativeMarketOrderDirections []*types.MatchedMarketDirection,
	stakingInfo *v2.FeeDiscountStakingInfo,
) ([]*v2.SpotBatchExecutionData, []*v2.DerivativeBatchExecutionData) {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.ExecuteMarketOrders")()

	spotResults := make([]*v2.SpotBatchExecutionData, len(spotMarketOrderIndicators))
	derivativeResults := make([]*v2.DerivativeBatchExecutionData, len(derivativeMarketOrderDirections))
	wg := new(sync.WaitGroup)

	// Spawn all spot market order goroutines
	// Legacy origin: exchange/abci.go EndBlocker loop over spotMarketOrderIndicators
	for idx, indicator := range spotMarketOrderIndicators {
		wg.Add(1)
		go func(i int, ind *v2.MarketOrderIndicator) {
			defer wg.Done()
			spotResults[i] = ba.executeSpotMarketOrder(ctx, ind, stakingInfo)
		}(idx, indicator)
	}

	// Spawn all derivative market order goroutines (concurrent with spot)
	// Legacy origin: exchange/abci.go EndBlocker loop over derivativeMarketOrderMarketDirections
	for idx, direction := range derivativeMarketOrderDirections {
		wg.Add(1)
		go func(i int, dir *types.MatchedMarketDirection) {
			defer wg.Done()
			derivativeResults[i] = ba.executeDerivativeMarketOrder(ctx, dir, stakingInfo)
		}(idx, direction)
	}

	// Wait for all market orders to complete
	wg.Wait()
	return spotResults, derivativeResults
}

// executeSpotMarketOrder executes market orders for a single spot market direction (buy or sell)
func (ba *BatchAuction) executeSpotMarketOrder(
	ctx sdk.Context,
	indicator *v2.MarketOrderIndicator,
	stakingInfo *v2.FeeDiscountStakingInfo,
) *v2.SpotBatchExecutionData {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.executeSpotMarketOrder")()

	marketID := common.HexToHash(indicator.MarketId)
	market := ba.spotKeeper.GetSpotMarket(ctx, marketID, true)
	if market == nil {
		return nil
	}

	executor := NewSpotMarketOrderExecutor(ctx, ba.spotKeeper, market, indicator.IsBuy)
	return executor.Execute(ctx, stakingInfo)
}

// executeDerivativeMarketOrder executes market orders for a single derivative market.
// This handles BOTH buy and sell market orders for the market in a single call
func (ba *BatchAuction) executeDerivativeMarketOrder(
	ctx sdk.Context,
	direction *types.MatchedMarketDirection,
	stakingInfo *v2.FeeDiscountStakingInfo,
) *v2.DerivativeBatchExecutionData {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.executeDerivativeMarketOrder")()

	marketID := direction.MarketId

	market, markPrice := ba.derivativeKeeper.GetDerivativeOrBinaryOptionsMarketWithMarkPrice(ctx, marketID, true)
	if market == nil {
		return nil
	}

	var funding *v2.PerpetualMarketFunding
	if market.GetIsPerpetual() {
		funding = ba.derivativeKeeper.GetPerpetualMarketFunding(ctx, marketID)
	}

	currentOpenNotional := ba.derivativeKeeper.GetOpenNotionalForMarket(ctx, marketID, markPrice)
	openNotionalCap := market.GetOpenNotionalCap()

	executor := NewDerivativeMarketOrderExecutor(
		ctx,
		ba.derivativeKeeper,
		market,
		markPrice,
		funding,
		currentOpenNotional,
		openNotionalCap,
	)
	return executor.Execute(ctx, stakingInfo)
}

// ExecuteLimitOrders executes the FBA (Frequent Batch Auction) for all limit orders
// across spot and derivative markets. This is the core batch auction algorithm where
// transient limit orders are matched against each other.
func (ba *BatchAuction) ExecuteLimitOrders(
	ctx sdk.Context,
	spotDirections []*types.MatchedMarketDirection,
	derivativeDirections []*types.MatchedMarketDirection,
	stakingInfo *v2.FeeDiscountStakingInfo,
	modifiedPositionCache v2.ModifiedPositionCache,
) ([]*v2.SpotBatchExecutionData, []*v2.DerivativeBatchExecutionData) {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.ExecuteLimitOrders")()

	spotResults := make([]*v2.SpotBatchExecutionData, len(spotDirections))
	derivativeResults := make([]*v2.DerivativeBatchExecutionData, len(derivativeDirections))
	wg := new(sync.WaitGroup)

	// Spawn all spot market goroutines
	for idx, direction := range spotDirections {
		wg.Add(1)
		go func(i int, marketID common.Hash) {
			defer wg.Done()
			spotResults[i] = ba.executeSpotLimitOrders(ctx, marketID, stakingInfo)
		}(idx, direction.MarketId)
	}

	// Spawn all derivative market goroutines (concurrent with spot)
	for idx, direction := range derivativeDirections {
		wg.Add(1)
		go func(i int, marketID common.Hash) {
			defer wg.Done()
			derivativeResults[i] = ba.executeDerivativeLimitOrders(ctx, marketID, stakingInfo, modifiedPositionCache)
		}(idx, direction.MarketId)
	}

	// Wait for all markets (spot and derivative) to complete
	wg.Wait()
	return spotResults, derivativeResults
}

// executeSpotLimitOrders executes the FBA for a single spot market.
func (ba *BatchAuction) executeSpotLimitOrders(
	ctx sdk.Context,
	marketID common.Hash,
	stakingInfo *v2.FeeDiscountStakingInfo,
) *v2.SpotBatchExecutionData {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.executeSpotLimitOrders")()

	market := ba.spotKeeper.GetSpotMarket(ctx, marketID, true)
	if market == nil {
		return nil
	}

	executor := NewSpotBatchAuctionExecutor(ctx, ba.spotKeeper, market)
	return executor.Execute(ctx, stakingInfo)
}

// executeDerivativeLimitOrders executes the FBA for a single derivative market.
func (ba *BatchAuction) executeDerivativeLimitOrders(
	ctx sdk.Context,
	marketID common.Hash,
	stakingInfo *v2.FeeDiscountStakingInfo,
	modifiedPositionCache v2.ModifiedPositionCache,
) *v2.DerivativeBatchExecutionData {
	defer ba.spotKeeper.Meter(ctx).FuncTiming(&ctx, "BatchAuction.executeDerivativeLimitOrders")()

	market, markPrice := ba.derivativeKeeper.GetDerivativeOrBinaryOptionsMarketWithMarkPrice(ctx, marketID, true)
	if market == nil {
		return nil
	}

	var funding *v2.PerpetualMarketFunding
	if market.GetIsPerpetual() {
		funding = ba.derivativeKeeper.GetPerpetualMarketFunding(ctx, marketID)
	}

	currentOpenNotional := ba.derivativeKeeper.GetOpenNotionalForMarket(ctx, marketID, markPrice)
	openNotionalCap := market.GetOpenNotionalCap()

	executor := NewDerivativeBatchAuctionExecutor(
		ctx,
		ba.derivativeKeeper,
		market,
		markPrice,
		funding,
		openNotionalCap,
		currentOpenNotional,
		modifiedPositionCache,
	)
	return executor.Execute(ctx, stakingInfo)
}
