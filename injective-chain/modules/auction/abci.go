package auction

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctiontypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

// subtractReserved returns a copy of coin with the amount reduced by the reserved amount for
// its denom. The result is clamped to zero to guard against any accounting inconsistency.
func subtractReserved(coin sdk.Coin, reserved map[string]math.Int) sdk.Coin {
	r, ok := reserved[coin.Denom]
	if !ok || !r.IsPositive() {
		return coin
	}
	coin.Amount = math.MaxInt(math.ZeroInt(), coin.Amount.Sub(r))
	return coin
}

// computeExchangeWithdrawalInjCap returns how much INJ the exchange may still send toward maxInjCap.
// totalInjInModule is the raw auction module INJ balance; reservedInj is the sum of outstanding
// voucher amounts for INJ (not auctionable). Without subtracting reservedInj, the cap would be
// understated and fresh INJ would be starved.
func computeExchangeWithdrawalInjCap(maxInjCap, totalInjInModule, reservedInj math.Int) math.Int {
	availableInj := math.MaxInt(math.ZeroInt(), totalInjInModule.Sub(reservedInj))
	return math.MaxInt(math.ZeroInt(), maxInjCap.Sub(availableInj))
}

func (am AppModule) burnBidAmountInModule(ctx sdk.Context, lastBidAmount math.Int) {
	defer am.keeper.Meter(ctx).FuncTiming(&ctx, "AppModule.burnBidAmountInModule")()

	injBurnAmount := sdk.NewCoin(chaintypes.InjectiveCoin, lastBidAmount)
	if err := am.bankKeeper.BurnCoins(ctx, auctiontypes.ModuleName, sdk.NewCoins(injBurnAmount)); err != nil {
		am.keeper.Logger(ctx).Error("failed to burn bid amount in auction module", "error", err)
	}
}

// sendBasketToWinner sends all positive balances from the auction module to the winner, capping INJ at maxInjCap.
// Coins that back outstanding vouchers from prior rounds are excluded so they cannot be double-spent.
// If an individual coin transfer fails (e.g. due to a send restriction), a voucher is created for the winner
// so they can claim the amount at a future time. The coins remain in the module account.
//
// Error-handling rationale (EndBlocker constraint):
//
// EndBlocker has no return value. Panicking halts the whole chain; returning early mid-settlement
// creates a permanently stuck auction (timer has elapsed but round never advances). Therefore all
// errors inside this method are handled with best-effort recovery:
//
//   - Voucher reservation read failure: we proceed with an empty map so the winner still receives
//     their basket. A KV-store read failure here indicates catastrophic state corruption that no
//     application-level guard can remedy. The marginal double-spend risk is accepted as the lesser
//     of two evils over denying the winner their rightful payout.
//
//   - Per-coin send failure: a voucher is created so the winner can claim the amount later.
//     The coins remain in the module account, so no value is lost.
//
//   - Voucher write failure after a send failure: this is a double-failure on corrupted state.
//     The error is logged; the coins remain in the module and are not permanently lost, but the
//     accounting record cannot be recovered without operator intervention.
func (am AppModule) sendBasketToWinner(ctx sdk.Context, auctionModuleAddress sdk.AccAddress, maxInjCap math.Int, winner sdk.AccAddress) {
	defer am.keeper.Meter(ctx).FuncTiming(&ctx, "AppModule.sendBasketToWinner")()

	reserved, err := am.keeper.GetVoucherReservedPerDenom(ctx)
	if err != nil {
		// A store read failure here indicates state corruption; the double-spend risk is accepted
		// as the lesser of two evils over leaving the winner with nothing.
		am.keeper.Logger(ctx).Error("failed to load voucher reservations for basket send, proceeding without deduction", "error", err)
		reserved = map[string]math.Int{}
	}

	coins := am.bankKeeper.GetAllBalances(ctx, auctionModuleAddress)
	for _, coin := range coins {
		coin = subtractReserved(coin, reserved)
		if coin.Denom == chaintypes.InjectiveCoin && coin.Amount.GT(maxInjCap) {
			coin.Amount = maxInjCap
		}
		if coin.Amount.IsPositive() {
			if err := am.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, winner, sdk.NewCoins(coin)); err != nil {
				am.keeper.Logger(ctx).Error("basket coin delivery to winner failed, creating voucher", "coin", coin.String(), "winner", winner.String(), "error", err)
				am.keeper.CreateVoucherOnFailedSend(ctx, winner, coin)
			}
		}
	}
}

// settleFinishedAuctionRound burns the bid INJ, sends the auction basket to the winner (with INJ capped),
// stores and emits the result, and clears the bid. No-op if there is no valid bid to settle.
// Returns an error if there was a bid but settlement failed (e.g. malformed bidder); the caller must not advance the round.
//
// Intentional design: sendBasketToWinner errors are not propagated. Bid INJ is burned before the basket
// is sent; if sendBasketToWinner returned an error and we aborted here, the INJ would be permanently gone
// but the bid record would not be deleted, creating an unrecoverable inconsistency on the next block.
// Individual delivery failures are instead handled inside sendBasketToWinner via vouchers (see its doc).
func (am AppModule) settleFinishedAuctionRound(ctx sdk.Context, auctionModuleAddress sdk.AccAddress, maxInjCap math.Int) error {
	defer am.keeper.Meter(ctx).FuncTiming(&ctx, "AppModule.settleFinishedAuctionRound")()

	lastBid := am.keeper.GetHighestBid(ctx)
	if lastBid == nil || !lastBid.Amount.Amount.IsPositive() || lastBid.Bidder == "" {
		return nil
	}
	lastBidder, err := sdk.AccAddressFromBech32(lastBid.Bidder)
	if err != nil {
		return err
	}

	am.burnBidAmountInModule(ctx, lastBid.Amount.Amount)
	am.sendBasketToWinner(ctx, auctionModuleAddress, maxInjCap, lastBidder)

	auctionRound := am.keeper.GetAuctionRound(ctx)
	am.keeper.SetLastAuctionResult(ctx, auctiontypes.LastAuctionResult{
		Winner: lastBid.Bidder,
		Amount: lastBid.Amount,
		Round:  auctionRound,
	})
	// nolint:errcheck //ignored on purpose
	ctx.EventManager().EmitTypedEvent(&auctiontypes.EventAuctionResult{
		Winner: lastBid.Bidder,
		Amount: lastBid.Amount,
		Round:  auctionRound,
	})
	am.keeper.DeleteBid(ctx)
	return nil
}

func (am AppModule) EndBlocker(ctx sdk.Context) {
	defer am.keeper.Meter(ctx).FuncTiming(&ctx, "EndBlocker")()

	// trigger auction settlement
	endingTimeStamp := am.keeper.GetEndingTimeStamp(ctx)

	if ctx.BlockTime().Unix() < endingTimeStamp {
		return
	}

	logger := ctx.Logger().With("module", "auction", "EndBlocker", ctx.BlockHeight())
	logger.Info("Settling auction round...", "blockTimestamp", ctx.BlockTime().Unix(), "endingTimeStamp", endingTimeStamp)
	auctionModuleAddress := am.accountKeeper.GetModuleAddress(auctiontypes.ModuleName)
	maxInjCap := am.keeper.GetParams(ctx).InjBasketMaxCap
	if err := am.settleFinishedAuctionRound(ctx, auctionModuleAddress, maxInjCap); err != nil {
		logger.Error("failed to settle auction round", "error", err)
		return
	}

	// advance auctionRound, endingTimestamp
	nextRound := am.keeper.AdvanceNextAuctionRound(ctx)
	nextEndingTimestamp := am.keeper.AdvanceNextEndingTimeStamp(ctx)

	// Single read: voucher map is unchanged until the next round's settlement (sweep / exchange
	// withdrawal do not create vouchers). Reuse for INJ cap and EventAuctionStart basket.
	// On error we proceed with an empty map. Returning early here would leave the auction stuck:
	// the round timer has already elapsed and settlement completed, but the round counter and
	// ending timestamp would not advance, so the EndBlocker would repeatedly fire without progress.
	voucherReserved, err := am.keeper.GetVoucherReservedPerDenom(ctx)
	if err != nil {
		logger.Error("failed to load voucher reservations", "error", err)
		voucherReserved = map[string]math.Int{}
	}

	// sweep ante fees from fees subaccount into module account (no INJ cap)
	am.keeper.SweepFeesSubaccountToModule(ctx)
	injInModule := am.bankKeeper.GetBalance(ctx, auctionModuleAddress, chaintypes.InjectiveCoin).Amount
	reservedInj := math.ZeroInt()
	if r, ok := voucherReserved[chaintypes.InjectiveCoin]; ok {
		reservedInj = r
	}
	exchangeWithdrawalInjCap := computeExchangeWithdrawalInjCap(maxInjCap, injInModule, reservedInj)
	// ping exchange module to flush fee for next round
	balances := am.exchangeKeeper.WithdrawAllAuctionBalances(ctx, exchangeWithdrawalInjCap)

	rawBasket := am.bankKeeper.GetAllBalances(ctx, auctionModuleAddress)
	newBasket := make(sdk.Coins, 0, len(rawBasket))
	for _, coin := range rawBasket {
		coin = subtractReserved(coin, voucherReserved)
		if coin.IsPositive() {
			newBasket = append(newBasket, coin)
		}
	}

	// for correctness, emit the correct INJ value in the new basket in the event the INJ balances exceed the cap
	newInjAmount := newBasket.AmountOf(chaintypes.InjectiveCoin)
	if newInjAmount.GT(maxInjCap) {
		excessInj := newInjAmount.Sub(maxInjCap)
		newBasket = newBasket.Sub(sdk.NewCoin(chaintypes.InjectiveCoin, excessInj))
	}

	// nolint:errcheck //ignored on purpose
	ctx.EventManager().EmitTypedEvent(&auctiontypes.EventAuctionStart{
		Round:           nextRound,
		EndingTimestamp: nextEndingTimestamp,
		NewBasket:       newBasket,
	})

	if len(balances) == 0 {
		logger.Info("😢 Received empty coin basket from exchange")
	} else {
		logger.Info("💰 Auction module received", balances.String(), "new auction basket is now", newBasket.String())
	}
}
