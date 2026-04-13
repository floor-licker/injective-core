package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	vouchertypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
)

// GetVoucherForAddress returns the outstanding voucher coin for the given denom and address.
// Returns a zero-amount coin when no voucher exists.
func (k *Keeper) GetVoucherForAddress(ctx sdk.Context, denom string, addr sdk.AccAddress) (coin sdk.Coin, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetVoucherForAddress")(&err)
	return k.vouchersAssistant.GetVoucher(ctx, denom, addr)
}

// SetVoucher persists an exact voucher amount, intended for genesis restore only.
// It overwrites any existing voucher for that (addr, denom). For incremental credits use
// vouchersAssistant.AddVoucher via module-specific helpers, not SetVoucher.
func (k *Keeper) SetVoucher(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	return k.vouchersAssistant.SetVoucher(ctx, addr, coin)
}

// GetAllVouchers returns every outstanding voucher across all denoms and addresses.
func (k *Keeper) GetAllVouchers(ctx sdk.Context) ([]vouchertypes.AddressVoucher, error) {
	return k.vouchersAssistant.GetAllVouchers(ctx)
}
