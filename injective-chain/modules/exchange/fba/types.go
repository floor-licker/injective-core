package fba

import (
	"cosmossdk.io/math"
)

// MatchResult captures the outcome of orderbook matching
type MatchResult struct {
	ClearingPrice    math.LegacyDec
	ClearingQuantity math.LegacyDec
	LastBuyPrice     math.LegacyDec
	LastSellPrice    math.LegacyDec
}
