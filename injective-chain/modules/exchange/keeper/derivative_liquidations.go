package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

type LiquidationMode int

const (
	LiquidationModeRegular LiquidationMode = iota
	LiquidationModeOffsetting
	LiquidationModeEmergencySettle
)

func getLiquidatorRewardShareRate(
	params v2.Params,
	//revive:disable:flag-parameter
	hasLiquidatorProvidedOrder bool,
	isWhiteKnightLiquidator bool,
) math.LegacyDec {
	if hasLiquidatorProvidedOrder && isWhiteKnightLiquidator {
		return params.WhiteKnightLiquidatorRewardShareRate
	}

	return params.LiquidatorRewardShareRate
}

func (k DerivativesMsgServer) handlePositiveLiquidationPayout(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	surplusAmount math.LegacyDec,
	liquidatorAddr sdk.AccAddress,
	positionSubaccountID common.Hash,
	liquidatorRewardShareRate math.LegacyDec,
) error {
	defer k.Meter(ctx).FuncTiming(&ctx, "handlePositiveLiquidationPayout")()

	insuranceFundOrAuctionPaymentAmount := surplusAmount.Mul(math.LegacyOneDec().Sub(liquidatorRewardShareRate)).TruncateInt()
	liquidatorPayout := surplusAmount.Sub(insuranceFundOrAuctionPaymentAmount.ToLegacyDec())

	if liquidatorPayout.IsPositive() {
		k.IncrementDepositOrSendToBank(ctx, types.SdkAddressToSubaccountID(liquidatorAddr), market.QuoteDenom, liquidatorPayout)
	}

	k.UpdateDepositWithDelta(ctx, positionSubaccountID, market.QuoteDenom, &types.DepositDelta{
		AvailableBalanceDelta: surplusAmount.Neg(),
		TotalBalanceDelta:     surplusAmount.Neg(),
	})

	if !insuranceFundOrAuctionPaymentAmount.IsPositive() {
		return nil
	}

	return k.MoveCoinsIntoInsuranceFund(ctx, market, insuranceFundOrAuctionPaymentAmount)
}

func (k DerivativesMsgServer) handlePositiveOffsettingLiquidationPayout(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	surplusAmount math.LegacyDec,
	liquidatorAddr sdk.AccAddress,
	liquidatorRewardShareRate math.LegacyDec,
	depositDeltas types.DepositDeltas,
) error {
	defer k.Meter(ctx).FuncTiming(&ctx, "handlePositiveOffsettingLiquidationPayout")()

	insuranceFundOrAuctionPaymentAmount := surplusAmount.Mul(math.LegacyOneDec().Sub(liquidatorRewardShareRate)).TruncateInt()
	liquidatorPayout := surplusAmount.Sub(insuranceFundOrAuctionPaymentAmount.ToLegacyDec())

	if liquidatorPayout.IsPositive() {
		depositDeltas.ApplyUniformDelta(types.SdkAddressToSubaccountID(liquidatorAddr), liquidatorPayout)
	}

	if !insuranceFundOrAuctionPaymentAmount.IsPositive() {
		return nil
	}

	return k.MoveCoinsIntoInsuranceFund(ctx, market, insuranceFundOrAuctionPaymentAmount)
}

// Four levels of escalation to retrieve the funds:
// 1: From trader's available balance
// 2: From trader's locked balance by cancelling his vanilla limit orders
// 3: From the insurance fund
// 4: Not enough funds available. Pause the market and socialize losses.
func (k DerivativesMsgServer) handleNegativeLiquidationPayout(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	positionSubaccountID common.Hash,
	lostFundsFromAvailableDuringPayout math.LegacyDec,
	isAllowingInsuranceFund bool,
) (shouldSettleMarket bool, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleNegativeLiquidationPayout")()

	shouldSettleMarket = false

	marketID := market.MarketID()
	liquidatedTraderDeposits := k.GetDeposit(ctx, positionSubaccountID, market.QuoteDenom)

	// defensive programming, orders should have been cancelled before this point
	if liquidatedTraderDeposits.HasTransientOrRestingVanillaLimitOrders() {
		k.CancelAllOrdersFromTraderInCurrentMarket(ctx, market, positionSubaccountID)
		k.CancelAllConditionalDerivativeOrdersBySubaccountIDAndMarket(ctx, market, positionSubaccountID)
	}

	availableBalanceAfterCancels := k.GetDeposit(ctx, positionSubaccountID, market.QuoteDenom).AvailableBalance
	retrievedFromCancellingOrders := availableBalanceAfterCancels.Sub(liquidatedTraderDeposits.AvailableBalance)
	lostFundsFromOrderCancels := retrievedFromCancellingOrders.Sub(math.LegacyMaxDec(math.LegacyZeroDec(), availableBalanceAfterCancels))

	k.EmitEvent(ctx, &v2.EventLostFundsFromLiquidation{
		MarketId:                           marketID.Hex(),
		SubaccountId:                       positionSubaccountID.Bytes(),
		LostFundsFromAvailableDuringPayout: lostFundsFromAvailableDuringPayout,
		LostFundsFromOrderCancels:          lostFundsFromOrderCancels,
	})

	k.IncrementMarketBalance(ctx, marketID, lostFundsFromAvailableDuringPayout.Add(lostFundsFromOrderCancels))

	if !availableBalanceAfterCancels.IsNegative() {
		return shouldSettleMarket, nil
	}

	absoluteDeficitAmount := availableBalanceAfterCancels.Abs()

	// trader has negative available balance, add the deficit amount to his position, because the negative balance is afterwards paid
	// by the insurance fund and through socialized loss during market settlement
	deposits := k.GetDeposit(ctx, positionSubaccountID, market.QuoteDenom)
	deposits.AvailableBalance = deposits.AvailableBalance.Add(absoluteDeficitAmount)
	deposits.TotalBalance = deposits.TotalBalance.Add(absoluteDeficitAmount)
	k.SetDeposit(ctx, positionSubaccountID, market.QuoteDenom, deposits)

	if !isAllowingInsuranceFund {
		shouldSettleMarket = true
		return shouldSettleMarket, nil
	}

	if absoluteDeficitAmount, err = k.PayDeficitFromInsuranceFund(ctx, marketID, absoluteDeficitAmount); err != nil {
		return shouldSettleMarket, err
	}

	if !absoluteDeficitAmount.IsPositive() {
		return shouldSettleMarket, nil
	}

	shouldSettleMarket = true
	return shouldSettleMarket, nil
}

func (k DerivativesMsgServer) EmergencySettleMarket(
	c context.Context, msg *v2.MsgEmergencySettleMarket,
) (*v2.MsgEmergencySettleMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "EmergencySettleMarket")()

	if !k.IsAdmin(ctx, msg.Sender) {
		return nil, sdkerrors.ErrUnauthorized
	}

	liquidatorAddr, _ := sdk.AccAddressFromBech32(msg.Sender)
	_, err := k.liquidatePosition(
		ctx,
		liquidatorAddr,
		common.HexToHash(msg.SubaccountId),
		common.HexToHash(msg.MarketId),
		nil,
		LiquidationModeEmergencySettle,
	)

	return &v2.MsgEmergencySettleMarketResponse{}, err
}

func (k DerivativesMsgServer) OffsetPosition(
	c context.Context, msg *v2.MsgOffsetPosition,
) (*v2.MsgOffsetPositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "OffsetPosition")()

	if !k.IsAdmin(ctx, msg.Sender) {
		return nil, sdkerrors.ErrUnauthorized
	}

	liquidatorAddr, _ := sdk.AccAddressFromBech32(msg.Sender)
	_, err := k.liquidatePosition(
		ctx,
		liquidatorAddr,
		common.HexToHash(msg.SubaccountId),
		common.HexToHash(msg.MarketId),
		nil,
		LiquidationModeOffsetting,
		msg.OffsettingSubaccountIds...,
	)

	return &v2.MsgOffsetPositionResponse{}, err
}

func (k DerivativesMsgServer) LiquidatePosition(
	c context.Context, msg *v2.MsgLiquidatePosition,
) (*v2.MsgLiquidatePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "LiquidatePosition")()

	liquidatorAddr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return k.liquidatePosition(
		ctx,
		liquidatorAddr,
		common.HexToHash(msg.SubaccountId),
		common.HexToHash(msg.MarketId),
		msg.Order,
		LiquidationModeRegular,
	)
}

func (k DerivativesMsgServer) prepareLiquidationMarketOrder(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	markPrice math.LegacyDec,
	funding *v2.PerpetualMarketFunding,
	position *v2.Position,
	positionSubaccountID common.Hash,
	liquidatorAddr sdk.AccAddress,
) (*v2.DerivativeMarketOrder, error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "prepareLiquidationMarketOrder")()

	marketOrderWorstPrice := position.GetLiquidationMarketOrderWorstPrice(markPrice, funding)

	liquidationMarketOrder := v2.NewMarketOrderForLiquidation(position, positionSubaccountID, liquidatorAddr, *marketOrderWorstPrice)

	subaccountNonce := k.IncrementSubaccountTradeNonce(ctx, positionSubaccountID)
	orderHash, err := liquidationMarketOrder.ComputeOrderHash(subaccountNonce.Nonce, market.MarketId)
	if err != nil {
		return nil, err
	}

	liquidationMarketOrder.OrderHash = orderHash.Bytes()

	return liquidationMarketOrder, nil
}

func (k DerivativesMsgServer) prepareLiquidatorOrder(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	markPrice math.LegacyDec,
	liquidatorOrder *v2.DerivativeOrder,
	liquidatorAddr sdk.AccAddress,
	liquidationMode LiquidationMode,
) (common.Hash, error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "prepareLiquidatorOrder")()

	liquidatorSubaccountID := types.MustGetSubaccountIDOrDeriveFromNonce(liquidatorAddr, liquidatorOrder.OrderInfo.SubaccountId)
	liquidatorOrder.OrderInfo.SubaccountId = liquidatorSubaccountID.Hex()
	metadata := k.GetSubaccountOrderbookMetadata(ctx, market.MarketID(), liquidatorSubaccountID, liquidatorOrder.IsBuy())

	isMaker := true
	liquidatorOrderHash, err := k.EnsureValidDerivativeOrder(ctx, liquidatorOrder, market, metadata, markPrice, false, nil, isMaker)

	// for emergency settling markets, we allow an invalid order, all order state changes are reverted later anyways
	if err != nil && liquidationMode != LiquidationModeEmergencySettle {
		return common.Hash{}, err
	}

	order := v2.NewDerivativeLimitOrder(liquidatorOrder, liquidatorAddr, liquidatorOrderHash)
	k.SetNewDerivativeLimitOrderWithMetadata(ctx, order, metadata, market.MarketID())

	return liquidatorOrderHash, nil
}

func (k DerivativesMsgServer) handleLiquidatorOrderPostExecution(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	marketID common.Hash,
	liquidatorOrder *v2.DerivativeOrder,
	liquidatorOrderHash common.Hash,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleLiquidatorOrderPostExecution")()

	isBuy := liquidatorOrder.IsBuy()
	subaccountID := liquidatorOrder.SubaccountID()
	orderAfterLiquidation := k.GetDerivativeLimitOrderBySubaccountIDAndHash(ctx, marketID, &isBuy, subaccountID, liquidatorOrderHash)

	if orderAfterLiquidation == nil || orderAfterLiquidation.Fillable.IsZero() {
		return
	}

	if err := k.CancelRestingDerivativeLimitOrder(
		ctx, market, orderAfterLiquidation.SubaccountID(), &isBuy, liquidatorOrderHash, true, true,
	); err != nil {
		k.Logger(ctx).Info(
			"CancelRestingDerivativeLimitOrder failed during LiquidatePosition of subaccount",
			"subaccountID", subaccountID.Hex(),
			"order", liquidatorOrder.String(),
			"err", err,
		)
		k.EmitEvent(
			ctx,
			v2.NewEventOrderCancelFail(
				marketID, subaccountID, orderAfterLiquidation.Hash().Hex(), orderAfterLiquidation.Cid(), err,
			),
		)
	}
}

func calculatePayout(
	fundsBeforeLiquidation math.LegacyDec,
	fundsAfterLiquidation math.LegacyDec,
) math.LegacyDec {
	if fundsBeforeLiquidation.IsNegative() {
		// if funds before liquidation are negative, then the initial negative balance should be included in the payout
		return fundsAfterLiquidation
	}
	return fundsAfterLiquidation.Sub(fundsBeforeLiquidation)
}

func calculateLostFundsFromAvailable(
	payout math.LegacyDec,
	//revive:disable:flag-parameter
	isMissingFunds bool,
	availableBalanceBeforeLiquidation math.LegacyDec,
) math.LegacyDec {
	if isMissingFunds {
		// balance is now negative, so trader lost all his available balance from liquidation
		return availableBalanceBeforeLiquidation
	} else if payout.IsNegative() {
		// balance is still positive, but negative payout still means trader lost some available balance from liquidation
		return payout.Abs()
	}
	return math.LegacyZeroDec()
}

func getOffsettingSettlementPrice(
	position *v2.Position,
	markPrice math.LegacyDec,
	funding *v2.PerpetualMarketFunding,
) (settlementPrice math.LegacyDec, isBankrupt bool) {
	bankruptcyPrice := position.GetBankruptcyPrice(funding)
	isBankrupt = (position.IsLong && markPrice.LTE(bankruptcyPrice)) || (position.IsShort() && markPrice.GTE(bankruptcyPrice))
	if isBankrupt {
		return bankruptcyPrice, true
	}

	return markPrice, false
}

func shouldHandlePositiveOffsettingLiquidationPayout(payout math.LegacyDec) (bool, error) {
	// defensive programming check
	if payout.IsNegative() {
		return false, errors.Wrapf(
			types.ErrPositionNotOffsettable,
			"non-bankrupt offsetting liquidation payout must be non-negative: %s",
			payout.String(),
		)
	}

	return payout.IsPositive(), nil
}

func parseSubaccountIDHashes(offsettingSubaccountIDs []string) []common.Hash {
	hashes := make([]common.Hash, 0, len(offsettingSubaccountIDs))
	for _, idStr := range offsettingSubaccountIDs {
		hashes = append(hashes, common.HexToHash(idStr))
	}
	return hashes
}

type offsetProcessResult struct {
	buyTrades          []*v2.DerivativeTradeLog
	sellTrades         []*v2.DerivativeTradeLog
	depositDeltas      types.DepositDeltas
	marketBalanceDelta math.LegacyDec
	remainingQuantity  math.LegacyDec
}

func (k DerivativesMsgServer) processOffsettingSubaccounts(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	settlementPrice math.LegacyDec,
	funding *v2.PerpetualMarketFunding,
	position *v2.Position,
	offsetIDs []common.Hash,
) (offsetProcessResult, error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "processOffsettingSubaccounts")()

	marketID := market.MarketID()
	remaining := position.Quantity

	res := offsetProcessResult{
		buyTrades:          []*v2.DerivativeTradeLog{},
		sellTrades:         []*v2.DerivativeTradeLog{},
		depositDeltas:      types.NewDepositDeltas(),
		marketBalanceDelta: math.LegacyZeroDec(),
	}

	for _, id := range offsetIDs {
		if remaining.IsZero() {
			break
		}

		pos := k.GetPosition(ctx, marketID, id)
		if pos == nil || pos.Quantity.IsZero() {
			continue
		}
		if pos.IsLong == position.IsLong {
			return offsetProcessResult{}, errors.Wrapf(types.ErrPositionNotOffsettable,
				"cannot offset same‑direction position %s in market %s", id.Hex(), marketID.Hex())
		}

		offsettingPosition := pos.Copy()
		offsettingPosition.ApplyFunding(funding)

		qty := math.LegacyMinDec(remaining, offsettingPosition.Quantity)
		if !qty.IsPositive() {
			continue
		}

		delta := &v2.PositionDelta{
			IsLong:            !offsettingPosition.IsLong,
			ExecutionQuantity: qty,
			ExecutionMargin:   math.LegacyZeroDec(),
			ExecutionPrice:    settlementPrice,
		}
		payout, _, _, pnl := offsettingPosition.ApplyPositionDelta(delta, math.LegacyZeroDec())
		if payout.IsNegative() {
			continue
		}

		k.CancelAllRestingDerivativeLimitOrdersForSubaccount(ctx, market, id, true, true)

		remaining = remaining.Sub(qty)

		chainPayout := market.NotionalToChainFormat(payout)
		res.marketBalanceDelta = res.marketBalanceDelta.Add(chainPayout.Neg())
		res.depositDeltas.ApplyUniformDelta(id, chainPayout)

		log := &v2.DerivativeTradeLog{
			SubaccountId:        id.Bytes(),
			PositionDelta:       delta,
			Payout:              payout,
			Fee:                 math.LegacyZeroDec(),
			OrderHash:           common.Hash{}.Bytes(),
			FeeRecipientAddress: common.Address{}.Bytes(),
			Pnl:                 pnl,
		}
		if offsettingPosition.IsLong {
			res.sellTrades = append(res.sellTrades, log)
		} else {
			res.buyTrades = append(res.buyTrades, log)
		}

		k.SavePosition(ctx, marketID, id, offsettingPosition)
	}

	res.remainingQuantity = remaining

	// Validate that at least some of the position was offset
	if remaining.Equal(position.Quantity) {
		offsetIDsStr := make([]string, len(offsetIDs))
		for i, id := range offsetIDs {
			offsetIDsStr[i] = id.Hex()
		}
		return offsetProcessResult{}, errors.Wrapf(types.ErrNoOffsettingPositionsFound,
			"no valid offsetting positions found from subaccounts [%v] in market %s", offsetIDsStr, marketID.Hex())
	}

	return res, nil
}

func (k DerivativesMsgServer) handleLiquidatedPosition(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	settlementPrice math.LegacyDec,
	funding *v2.PerpetualMarketFunding,
	position *v2.Position,
	positionSubaccountID common.Hash,
	liquidatorAddr sdk.AccAddress,
	liquidatorRewardShareRate math.LegacyDec,
	isBankrupt bool,
	res offsetProcessResult,
) error {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleLiquidatedPosition")()

	buyTrades, sellTrades, deltas, mktBalDelta, remainingQty :=
		res.buyTrades, res.sellTrades, res.depositDeltas, res.marketBalanceDelta, res.remainingQuantity

	wasLong := position.IsLong
	closingQuantity := position.Quantity.Sub(remainingQty)
	var (
		payout   math.LegacyDec
		pnl      math.LegacyDec
		liqDelta *v2.PositionDelta
	)
	if isBankrupt {
		pnl, liqDelta = position.ApplyBankruptCloseWithoutPayouts(settlementPrice, closingQuantity)
		payout = math.LegacyZeroDec()
	} else {
		liqDelta = &v2.PositionDelta{
			IsLong:            !position.IsLong,
			ExecutionQuantity: closingQuantity,
			ExecutionMargin:   math.LegacyZeroDec(),
			ExecutionPrice:    settlementPrice,
		}
		payout, _, _, pnl = position.ApplyPositionDelta(liqDelta, math.LegacyZeroDec())
		shouldHandlePositivePayout, err := shouldHandlePositiveOffsettingLiquidationPayout(payout)
		if err != nil {
			return err
		}
		if shouldHandlePositivePayout {
			chainPayout := market.NotionalToChainFormat(payout)
			mktBalDelta = mktBalDelta.Add(chainPayout.Neg())
			if err := k.handlePositiveOffsettingLiquidationPayout(
				ctx,
				market,
				chainPayout,
				liquidatorAddr,
				liquidatorRewardShareRate,
				deltas,
			); err != nil {
				return err
			}
		}
	}

	trade := &v2.DerivativeTradeLog{
		SubaccountId: positionSubaccountID.Bytes(), PositionDelta: liqDelta, Payout: payout, Pnl: pnl,
		Fee: math.LegacyZeroDec(), OrderHash: common.Hash{}.Bytes(), FeeRecipientAddress: common.Address{}.Bytes(),
	}
	if wasLong {
		sellTrades = append(sellTrades, trade)
	} else {
		buyTrades = append(buyTrades, trade)
	}

	k.SavePosition(ctx, market.MarketID(), positionSubaccountID, position)
	k.SetMarketBalance(ctx, market.MarketID(), k.GetMarketBalance(ctx, market.MarketID()).Add(mktBalDelta))

	// OI tracks both sides (long + short), each reduced by the closed quantity
	closedQty := liqDelta.ExecutionQuantity
	openInterestDelta := closedQty.MulInt64(-2)
	k.ApplyOpenInterestDeltaForMarket(ctx, market.MarketID(), openInterestDelta)

	var cumulativeFunding math.LegacyDec
	if funding != nil {
		cumulativeFunding = funding.CumulativeFunding
	}
	batch := func(isBuy, isLiq bool, trades []*v2.DerivativeTradeLog) *v2.EventBatchDerivativeExecution {
		return &v2.EventBatchDerivativeExecution{MarketId: market.MarketID().String(), IsBuy: isBuy, IsLiquidation: isLiq,
			ExecutionType: v2.ExecutionType_OffsettingPosition, Trades: trades, CumulativeFunding: &cumulativeFunding}
	}
	k.EmitEvent(ctx, batch(true, !wasLong, buyTrades))
	k.EmitEvent(ctx, batch(false, wasLong, sellTrades))

	for _, id := range deltas.GetSortedSubaccountKeys() {
		k.UpdateDepositWithDeltaWithoutBankCharge(ctx, id, market.GetQuoteDenom(), deltas[id])
	}

	return nil
}

func (k DerivativesMsgServer) handleOffsettingPositions(
	ctx sdk.Context,
	market *v2.DerivativeMarket,
	markPrice math.LegacyDec,
	funding *v2.PerpetualMarketFunding,
	position *v2.Position,
	positionSubaccountID common.Hash,
	liquidatorAddr sdk.AccAddress,
	offsettingSubaccountIDs ...string,
) error {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleOffsettingPositions")()

	settlementPrice, isBankrupt := getOffsettingSettlementPrice(position, markPrice, funding)
	liquidatorRewardShareRate := k.GetCachedParams(ctx).WhiteKnightLiquidatorRewardShareRate
	offsettingSubaccountIDHashes := parseSubaccountIDHashes(offsettingSubaccountIDs)

	res, err := k.processOffsettingSubaccounts(
		ctx,
		market,
		settlementPrice,
		funding,
		position,
		offsettingSubaccountIDHashes,
	)
	if err != nil {
		return err
	}

	return k.handleLiquidatedPosition(
		ctx,
		market,
		settlementPrice,
		funding,
		position,
		positionSubaccountID,
		liquidatorAddr,
		liquidatorRewardShareRate,
		isBankrupt,
		res,
	)
}

func (k DerivativesMsgServer) liquidatePosition(
	c context.Context,
	liquidatorAddr sdk.AccAddress,
	liquidatedSubaccountID,
	marketID common.Hash,
	liquidatorOrder *v2.DerivativeOrder,
	liquidationMode LiquidationMode,
	offsettingSubaccountIDs ...string,
) (*v2.MsgLiquidatePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "liquidatePosition")()

	cacheCtx, writeCache := ctx.CacheContext()

	positionSubaccountID := liquidatedSubaccountID
	isOffsettingSubaccount := liquidationMode == LiquidationModeOffsetting
	isEmergencySettlingMarket := liquidationMode == LiquidationModeEmergencySettle

	// 1. Reject if derivative market id does not reference an active derivative market
	market, markPrice := k.GetDerivativeMarketWithMarkPrice(cacheCtx, marketID, true)
	if market == nil {
		k.Logger(ctx).Error("active derivative market doesn't exist", "marketID", marketID.Hex())

		return nil, errors.Wrapf(types.ErrDerivativeMarketNotFound, "active derivative market for marketID %s not found", marketID.Hex())
	}

	position := k.GetPosition(cacheCtx, marketID, positionSubaccountID)
	if position == nil || position.Quantity.IsZero() {

		return nil, errors.Wrapf(types.ErrPositionNotFound, "subaccountID %s marketID %s", positionSubaccountID.Hex(), marketID.Hex())
	}

	var funding *v2.PerpetualMarketFunding
	if market.IsPerpetual {
		funding = k.GetPerpetualMarketFunding(cacheCtx, marketID)
	}

	liquidationPrice := position.GetLiquidationPrice(market.MaintenanceMarginRatio, funding)
	shouldLiquidate := (position.IsLong && markPrice.LTE(liquidationPrice)) || (position.IsShort() && markPrice.GTE(liquidationPrice))

	if !shouldLiquidate {
		return nil, errors.Wrapf(
			types.ErrPositionNotLiquidable,
			"%s position liquidation price is %s but mark price is %s",
			position.GetDirectionString(),
			liquidationPrice.String(),
			markPrice.String(),
		)
	}

	// Step 1a: Cancel all limit orders created by the position holder in the given market
	k.CancelAllTransientDerivativeLimitOrdersBySubaccountID(cacheCtx, market, positionSubaccountID)
	k.CancelAllRestingDerivativeLimitOrdersForSubaccount(cacheCtx, market, positionSubaccountID, true, true)

	positionState := v2.ApplyFundingAndGetUpdatedPositionState(position, funding)
	k.SavePosition(cacheCtx, marketID, positionSubaccountID, positionState.Position)

	// Step 1b: Cancel all market orders created by the position holder in the given market
	k.CancelAllDerivativeMarketOrdersBySubaccountID(cacheCtx, market, positionSubaccountID, marketID)

	// Step 1c: Cancel all conditional orders created by the position holder in the given market
	k.CancelAllConditionalDerivativeOrdersBySubaccountIDAndMarket(cacheCtx, market, positionSubaccountID)

	if isOffsettingSubaccount {
		if err := k.handleOffsettingPositions(
			cacheCtx,
			market,
			markPrice,
			funding,
			position,
			positionSubaccountID,
			liquidatorAddr,
			offsettingSubaccountIDs...,
		); err != nil {
			return nil, err
		}

		writeCache()
		return &v2.MsgLiquidatePositionResponse{}, nil
	}

	liquidationMarketOrder, err := k.prepareLiquidationMarketOrder(
		cacheCtx,
		market,
		markPrice,
		funding,
		position,
		positionSubaccountID,
		liquidatorAddr,
	)
	if err != nil {
		return nil, err
	}

	liquidatorRewardShareRate := getLiquidatorRewardShareRate(
		k.GetCachedParams(ctx),
		liquidatorOrder != nil,
		k.IsWhiteKnightLiquidator(ctx, liquidatorAddr.String()),
	)

	if isEmergencySettlingMarket {
		var orderType v2.OrderType

		if position.IsLong {
			orderType = v2.OrderType_BUY
		} else {
			orderType = v2.OrderType_SELL
		}

		liquidatorOrder = &v2.DerivativeOrder{
			MarketId: marketID.Hex(),
			OrderInfo: v2.OrderInfo{
				SubaccountId: "0",
				Price:        markPrice,
				Quantity:     position.Quantity,
			},
			OrderType: orderType,
			Margin:    position.Quantity.Mul(markPrice),
		}
	}

	var liquidatorOrderHash common.Hash
	hasLiquidatorOrder := liquidatorOrder != nil

	if hasLiquidatorOrder {
		liquidatorOrderHash, err = k.prepareLiquidatorOrder(cacheCtx, market, markPrice, liquidatorOrder, liquidatorAddr, liquidationMode)
		if err != nil {
			return nil, err
		}
	}

	positionStates := v2.NewPositionStates()
	positionCache := make(map[common.Hash]*v2.Position)

	fundsBeforeLiquidation := k.GetSpendableFunds(cacheCtx, positionSubaccountID, market.QuoteDenom)
	availableBalanceBeforeLiquidation := k.GetDeposit(cacheCtx, positionSubaccountID, market.QuoteDenom).AvailableBalance

	_, isMarketSolvent, err := k.ExecuteDerivativeMarketOrderImmediately(
		cacheCtx, market, markPrice, funding, liquidationMarketOrder, positionStates, positionCache, true,
	)

	if err != nil {
		return nil, err
	}

	if !isMarketSolvent {
		if err := k.PauseMarketAndScheduleForSettlement(ctx, market.MarketID(), true); err != nil {
			return nil, err
		}
		return &v2.MsgLiquidatePositionResponse{}, nil
	}

	if hasLiquidatorOrder {
		k.handleLiquidatorOrderPostExecution(cacheCtx, market, marketID, liquidatorOrder, liquidatorOrderHash)
	}

	fundsAfterLiquidation := k.GetSpendableFunds(cacheCtx, positionSubaccountID, market.QuoteDenom)
	availableBalanceAfterLiquidation := k.GetDeposit(cacheCtx, positionSubaccountID, market.QuoteDenom).AvailableBalance

	payout := calculatePayout(fundsBeforeLiquidation, fundsAfterLiquidation)
	isMissingFunds := payout.IsNegative() && availableBalanceAfterLiquidation.IsNegative()

	shouldSettleMarketFromLiquidation := false
	lostFundsFromAvailableDuringPayout := calculateLostFundsFromAvailable(payout, isMissingFunds, availableBalanceBeforeLiquidation)

	// if payout is positive, then trader lost position margin + PNL which we cannot get here, but which is emitted as EventBatchDerivativeExecution
	if isMissingFunds {
		if shouldSettleMarketFromLiquidation, err = k.handleNegativeLiquidationPayout(
			cacheCtx,
			market,
			positionSubaccountID,
			lostFundsFromAvailableDuringPayout,
			!isOffsettingSubaccount,
		); err != nil {

			return nil, err
		}
	} else if payout.IsPositive() {
		surplusAmount := payout
		if err = k.handlePositiveLiquidationPayout(
			cacheCtx,
			market,
			surplusAmount,
			liquidatorAddr,
			positionSubaccountID,
			liquidatorRewardShareRate,
		); err != nil {
			return nil, err
		}
	}

	if !isMissingFunds {
		// if missing funds this event is already emitted inside handleNegativeLiquidationPayout
		k.EmitEvent(cacheCtx, &v2.EventLostFundsFromLiquidation{
			MarketId:                           marketID.Hex(),
			SubaccountId:                       positionSubaccountID.Bytes(),
			LostFundsFromAvailableDuringPayout: lostFundsFromAvailableDuringPayout,
			LostFundsFromOrderCancels:          math.LegacyZeroDec(),
		})

		k.IncrementMarketBalance(cacheCtx, marketID, lostFundsFromAvailableDuringPayout)
	}

	shouldSettleMarket := shouldSettleMarketFromLiquidation

	if isEmergencySettlingMarket && !shouldSettleMarket {
		return nil, types.ErrInvalidEmergencySettle
	}

	if shouldSettleMarket {
		if err = k.PauseMarketAndScheduleForSettlement(ctx, market.MarketID(), true); err != nil {
			return nil, err
		}
	} else {
		writeCache()
	}

	return &v2.MsgLiquidatePositionResponse{}, nil
}
