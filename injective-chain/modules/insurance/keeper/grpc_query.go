package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	vouchertypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/insurance/types"
)

var _ types.QueryServer = &Keeper{}

// InsuranceParams is grpc implementation to return module params
func (k *Keeper) InsuranceParams(c context.Context, _ *types.QueryInsuranceParamsRequest) (*types.QueryInsuranceParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "InsuranceParams")()

	params := k.GetParams(ctx)

	res := &types.QueryInsuranceParamsResponse{
		Params: params,
	}

	return res, nil
}

// InsuranceFund is grpc implementation to return the insurance fund for a given derivative market
func (k *Keeper) InsuranceFund(c context.Context, request *types.QueryInsuranceFundRequest) (*types.QueryInsuranceFundResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "InsuranceFund")()

	fund := k.GetInsuranceFund(ctx, common.HexToHash(request.MarketId))

	res := &types.QueryInsuranceFundResponse{
		Fund: fund,
	}

	return res, nil
}

// InsuranceFunds is grpc implementation to return all the insurance funds
func (k *Keeper) InsuranceFunds(c context.Context, request *types.QueryInsuranceFundsRequest) (*types.QueryInsuranceFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "InsuranceFunds")()

	funds := k.GetAllInsuranceFunds(ctx)

	res := &types.QueryInsuranceFundsResponse{
		Funds: funds,
	}

	return res, nil
}

// EstimatedRedemptions is grpc implementation to return estimated redemptions from user owned shared tokens
func (k *Keeper) EstimatedRedemptions(c context.Context, request *types.QueryEstimatedRedemptionsRequest) (*types.QueryEstimatedRedemptionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "EstimatedRedemptions")()

	address, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	res := &types.QueryEstimatedRedemptionsResponse{
		Amount: k.GetEstimatedRedemptions(ctx, address, common.HexToHash(request.MarketId)),
	}

	return res, nil
}

// PendingRedemptions is grpc implementation to return estimated pending redemption at the time of claim
func (k *Keeper) PendingRedemptions(c context.Context, request *types.QueryPendingRedemptionsRequest) (*types.QueryPendingRedemptionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "PendingRedemptions")()

	address, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	res := &types.QueryPendingRedemptionsResponse{
		Amount: k.GetPendingRedemptions(ctx, address, common.HexToHash(request.MarketId)),
	}

	return res, nil
}

func (k *Keeper) InsuranceModuleState(c context.Context, _ *types.QueryModuleStateRequest) (res *types.QueryModuleStateResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "InsuranceModuleState")(&err)

	var avs []vouchertypes.AddressVoucher

	avs, err = k.GetAllVouchers(ctx)
	if err != nil {
		return nil, err
	}

	res = &types.QueryModuleStateResponse{
		State: &types.GenesisState{
			Params:                         k.GetParams(ctx),
			InsuranceFunds:                 k.GetAllInsuranceFunds(ctx),
			RedemptionSchedule:             k.GetAllInsuranceFundRedemptions(ctx),
			NextShareDenomId:               k.ExportNextShareDenomId(ctx),
			NextRedemptionScheduleId:       k.ExportNextRedemptionScheduleId(ctx),
			FailedRedemptionSchedules:      k.GetAllFailedRedemptionSchedules(ctx),
			NextFailedRedemptionScheduleId: k.ExportNextFailedRedemptionScheduleId(ctx),
			Vouchers:                       avs,
		},
	}

	return res, nil
}

// FailedRedemptions returns all failed redemption schedules
func (k *Keeper) FailedRedemptions(c context.Context, _ *types.QueryFailedRedemptionsRequest) (*types.QueryFailedRedemptionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "FailedRedemptions")()
	schedules := k.GetAllFailedRedemptionSchedules(ctx)

	res := &types.QueryFailedRedemptionsResponse{
		Schedules: schedules,
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

	var addr sdk.AccAddress
	addr, err = sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	var voucher sdk.Coin
	voucher, err = k.GetVoucherForAddress(ctx, req.Denom, addr)
	if err != nil {
		return nil, err
	}

	res = &types.QueryVoucherResponse{Voucher: voucher}
	return res, nil
}
