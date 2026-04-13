package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
	vouchertypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) AuctionParams(c context.Context, _ *types.QueryAuctionParamsRequest) (*types.QueryAuctionParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "AuctionParams")()

	params := k.GetParams(ctx)

	res := &types.QueryAuctionParamsResponse{
		Params: params,
	}
	return res, nil
}

func (k *Keeper) CurrentAuctionBasket(c context.Context, _ *types.QueryCurrentAuctionBasketRequest) (*types.QueryCurrentAuctionBasketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "CurrentAuctionBasket")()

	auctionModuleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	coins := k.bankKeeper.GetAllBalances(ctx, auctionModuleAddress)
	lastBid := k.GetHighestBid(ctx)
	maxCap := k.GetParams(ctx).InjBasketMaxCap

	reserved, err := k.GetVoucherReservedPerDenom(ctx)
	if err != nil {
		reserved = map[string]math.Int{}
	}

	currentBasketCoins := make([]sdk.Coin, 0, len(coins))
	for _, coin := range coins {
		if r, ok := reserved[coin.Denom]; ok && r.IsPositive() {
			coin.Amount = math.MaxInt(math.ZeroInt(), coin.Amount.Sub(r))
		}
		if coin.Denom == chaintypes.InjectiveCoin {
			coin = coin.SubAmount(lastBid.Amount.Amount)
			if coin.Amount.IsNegative() {
				coin.Amount = math.ZeroInt()
			}
			if coin.Amount.GT(maxCap) {
				coin.Amount = maxCap
			}
		}
		if coin.Amount.IsPositive() {
			currentBasketCoins = append(currentBasketCoins, coin)
		}
	}

	closingTime := k.GetEndingTimeStamp(ctx)
	res := &types.QueryCurrentAuctionBasketResponse{
		AuctionRound:       k.GetAuctionRound(ctx),
		AuctionClosingTime: uint64(closingTime),
		HighestBidAmount:   lastBid.Amount.Amount,
		HighestBidder:      lastBid.Bidder,
		Amount:             sdk.NewCoins(currentBasketCoins...),
	}
	return res, nil
}

func (k *Keeper) AuctionModuleState(c context.Context, _ *types.QueryModuleStateRequest) (res *types.QueryModuleStateResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "AuctionModuleState")(&err)

	avs, err := k.GetAllVouchers(ctx)
	if err != nil {
		return nil, err
	}

	res = &types.QueryModuleStateResponse{
		State: &types.GenesisState{
			Params:                 k.GetParams(ctx),
			AuctionRound:           k.GetAuctionRound(ctx),
			HighestBid:             k.GetHighestBid(ctx),
			AuctionEndingTimestamp: k.GetEndingTimeStamp(ctx),
			LastAuctionResult:      k.GetLastAuctionResult(ctx),
			Vouchers:               avs,
		},
	}
	return res, nil
}

func (k *Keeper) LastAuctionResult(c context.Context, _ *types.QueryLastAuctionResultRequest) (*types.QueryLastAuctionResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "LastAuctionResult")()

	res := &types.QueryLastAuctionResultResponse{
		LastAuctionResult: k.GetLastAuctionResult(ctx),
	}
	return res, nil
}

func (k *Keeper) Vouchers(c context.Context, req *types.QueryVouchersRequest) (res *types.QueryVouchersResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "Vouchers")(&err)

	var avs []vouchertypes.AddressVoucher

	if req.Denom == "" {
		avs, err = k.vouchersAssistant.GetAllVouchers(ctx)
		if err != nil {
			return nil, err
		}
		res = &types.QueryVouchersResponse{Vouchers: avs}
		return res, nil
	}

	avs, err = k.vouchersAssistant.GetVouchersForDenom(ctx, req.Denom)
	if err != nil {
		return nil, err
	}
	res = &types.QueryVouchersResponse{Vouchers: avs}
	return res, nil
}

func (k *Keeper) Voucher(c context.Context, req *types.QueryVoucherRequest) (res *types.QueryVoucherResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "Voucher")(&err)

	if req.Denom == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "denom is required")
	}
	if req.Address == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "address is required")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	voucher, err := k.GetVoucherForAddress(ctx, req.Denom, addr)
	if err != nil {
		return nil, err
	}

	res = &types.QueryVoucherResponse{Voucher: voucher}
	return res, nil
}
