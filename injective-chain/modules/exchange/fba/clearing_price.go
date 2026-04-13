package fba

import (
	"cosmossdk.io/math"
)

// baseClearingStrategy provides shared clearing price calculation methods.
type baseClearingStrategy struct{}

// midpointPrice calculates the average of buy and sell prices.
func (baseClearingStrategy) midpointPrice(lastBuyPrice, lastSellPrice math.LegacyDec) math.LegacyDec {
	return lastBuyPrice.Add(lastSellPrice).Quo(math.LegacyNewDec(2))
}

// selectPriceByReference selects the clearing price based on a reference price.
// Returns lastBuyPrice if <= reference, lastSellPrice if >= reference, otherwise reference.
func (baseClearingStrategy) selectPriceByReference(lastBuyPrice, lastSellPrice, referencePrice math.LegacyDec) math.LegacyDec {
	if lastBuyPrice.LTE(referencePrice) {
		return lastBuyPrice
	}
	if lastSellPrice.GTE(referencePrice) {
		return lastSellPrice
	}
	return referencePrice
}

type SpotClearingPriceStrategy struct {
	baseClearingStrategy
	midMarketPrice *math.LegacyDec
}

func NewSpotClearingPriceStrategy(midMarketPrice *math.LegacyDec) *SpotClearingPriceStrategy {
	return &SpotClearingPriceStrategy{midMarketPrice: midMarketPrice}
}

func (s *SpotClearingPriceStrategy) Calculate(lastBuyPrice, lastSellPrice, clearingQuantity math.LegacyDec) math.LegacyDec {
	if !clearingQuantity.IsPositive() {
		return math.LegacyZeroDec()
	}

	if s.midMarketPrice != nil {
		return s.selectPriceByReference(lastBuyPrice, lastSellPrice, *s.midMarketPrice)
	}

	// Edge case when a resting orderbook does not exist, so no other choice
	return s.midpointPrice(lastBuyPrice, lastSellPrice)
}

type DerivativeClearingPriceStrategy struct {
	baseClearingStrategy
	markPrice      math.LegacyDec
	midMarketPrice *math.LegacyDec
}

func NewDerivativeClearingPriceStrategy(markPrice math.LegacyDec, midMarketPrice *math.LegacyDec) *DerivativeClearingPriceStrategy {
	return &DerivativeClearingPriceStrategy{
		markPrice:      markPrice,
		midMarketPrice: midMarketPrice,
	}
}

func (d *DerivativeClearingPriceStrategy) Calculate(lastBuyPrice, lastSellPrice, clearingQuantity math.LegacyDec) math.LegacyDec {
	if !clearingQuantity.IsPositive() {
		return math.LegacyZeroDec()
	}

	// Edge case: no resting orderbook and no mark price
	if d.midMarketPrice == nil && d.markPrice.IsNil() {
		return d.midpointPrice(lastBuyPrice, lastSellPrice)
	}

	// Oracle fallback (no mid-market price)
	if d.midMarketPrice == nil {
		return d.selectPriceByReference(lastBuyPrice, lastSellPrice, d.markPrice)
	}

	// Regular case with mid-market price
	return d.calculateRegular(lastBuyPrice, lastSellPrice)
}

func (d *DerivativeClearingPriceStrategy) calculateRegular(lastBuyPrice, lastSellPrice math.LegacyDec) math.LegacyDec {
	if lastBuyPrice.LTE(*d.midMarketPrice) {
		return lastBuyPrice
	}

	if lastSellPrice.GTE(*d.midMarketPrice) {
		return lastSellPrice
	}

	// Fall back to oracle price if available
	if !d.markPrice.IsNil() {
		return d.selectPriceByReference(lastBuyPrice, lastSellPrice, d.markPrice)
	}

	return *d.midMarketPrice
}
