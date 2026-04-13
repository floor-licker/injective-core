package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	vouchertypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
)

// GetVoucherForAddress returns the outstanding voucher coin for the given denom and address.
// Returns a zero-amount coin when no voucher exists.
func (k *Keeper) GetVoucherForAddress(ctx sdk.Context, denom string, addr sdk.AccAddress) (coin sdk.Coin, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetVoucherForAddress")(&err)
	coin, err = k.vouchersAssistant.GetVoucher(ctx, denom, addr)
	return coin, err
}

// SetVoucher persists an exact voucher amount, intended for genesis restore only.
// It overwrites any existing voucher for that (addr, denom). On incremental failure paths,
// use CreateVoucherOnFailedSend (which calls AddVoucher) instead.
func (k *Keeper) SetVoucher(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	return k.vouchersAssistant.SetVoucher(ctx, addr, coin)
}

// GetAllVouchers returns every outstanding voucher across all denoms and addresses.
func (k *Keeper) GetAllVouchers(ctx sdk.Context) ([]vouchertypes.AddressVoucher, error) {
	return k.vouchersAssistant.GetAllVouchers(ctx)
}

// GetVoucherReservedPerDenom returns the total outstanding voucher amount for each denom
// across all addresses. This is used to ensure that funds committed to outstanding vouchers
// are not re-included in subsequent auction basket sends.
func (k *Keeper) GetVoucherReservedPerDenom(ctx sdk.Context) (map[string]math.Int, error) {
	avs, err := k.GetAllVouchers(ctx)
	if err != nil {
		return nil, err
	}
	reserved := make(map[string]math.Int, len(avs))
	for _, av := range avs {
		existing, ok := reserved[av.Voucher.Denom]
		if !ok {
			existing = math.ZeroInt()
		}
		reserved[av.Voucher.Denom] = existing.Add(av.Voucher.Amount)
	}
	return reserved, nil
}

// CreateVoucherOnFailedSend records a voucher claim for amount on behalf of addr when a
// bank transfer fails. The coins remain in the module account so the user can later
// claim them via MsgClaimVoucher. Uses AddVoucher to merge with any existing voucher for the
// same (denom, addr) pair.
//
// The error from AddVoucher is intentionally not returned. This method is called from
// sendBasketToWinner, which runs inside EndBlocker. Propagating the error would surface in
// settleFinishedAuctionRound after bid INJ has already been burned — aborting at that point
// would leave the bid record intact while the INJ is gone, creating an unrecoverable state
// on the next block. A KV-store write failure here indicates catastrophic state corruption;
// the best achievable outcome is to log it so operators are alerted.
func (k *Keeper) CreateVoucherOnFailedSend(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) {
	if err := k.vouchersAssistant.AddVoucher(ctx, addr, amount); err != nil {
		k.Logger(ctx).Error(
			"failed to create voucher after basket delivery failure",
			"address", addr.String(),
			"coin", amount.String(),
			"error", err,
		)
	}
}
