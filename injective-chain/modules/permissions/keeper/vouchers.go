package keeper

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetVoucherForAddress returns the outstanding voucher coin for the given denom and address.
// Returns a zero-amount coin when no voucher exists.
func (k Keeper) GetVoucherForAddress(ctx sdk.Context, denom string, addr sdk.AccAddress) (coin sdk.Coin, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetVoucherForAddress")(&err)
	coin, err = k.vouchersAssistant.GetVoucher(ctx, denom, addr)
	return coin, err
}

// rerouteToVoucherOnFail intercepts a failed send inside consensus-critical code and
// accumulates the amount into a voucher for toAddr, redirecting the bank send to the
// module account so the block execution does not fail.
// If the context key is not set the original address and error are returned unchanged.
func (k Keeper) rerouteToVoucherOnFail(ctx sdk.Context, toAddr sdk.AccAddress, amount sdk.Coin, origErr error) (_ sdk.AccAddress, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "rerouteToVoucherOnFail")(&err)

	if ctx.Value(baseapp.DoNotFailFastSendContextKey) == nil {
		return toAddr, origErr
	}

	if err = k.vouchersAssistant.AddVoucher(ctx, toAddr, amount); err != nil {
		return toAddr, errors.Wrapf(err, "can't accumulate voucher for address, tried to reroute token send after error: %s", origErr.Error())
	}

	return authtypes.NewModuleAddress(k.ModuleName()), nil
}
