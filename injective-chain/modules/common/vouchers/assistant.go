package vouchers

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	vouchertypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

var delim = []byte("|")

// VouchersAssistant centralizes all voucher CRUD logic so that any module can
// support the voucher pattern by implementing VoucherKeeper and wiring this component.
//
//nolint:revive // The name VouchersAssistant is more readable than Assistant
type VouchersAssistant struct {
	keeper     VoucherKeeper
	bankKeeper BankKeeper
}

// NewVouchersAssistant creates a new VouchersAssistant backed by the given keeper and bank keeper.
func NewVouchersAssistant(keeper VoucherKeeper, bankKeeper BankKeeper) *VouchersAssistant {
	return &VouchersAssistant{
		keeper:     keeper,
		bankKeeper: bankKeeper,
	}
}

// GetVoucher returns the outstanding voucher coin for addr / denom.
// If no voucher exists it returns a zero-amount coin for that denom.
func (a *VouchersAssistant) GetVoucher(ctx sdk.Context, denom string, addr sdk.AccAddress) (coin sdk.Coin, err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "GetVoucher")(&err)

	if err = sdk.ValidateDenom(denom); err != nil {
		return sdk.Coin{}, errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: %s", denom)
	}

	store := a.keeper.GetVouchersStore(ctx)
	bz := store.Get(voucherKey(denom, addr))
	if len(bz) == 0 {
		return sdk.NewInt64Coin(denom, 0), nil
	}

	var amount math.Int
	if err := amount.Unmarshal(bz); err != nil {
		return sdk.NewInt64Coin(denom, 0), err
	}
	return sdk.NewCoin(denom, amount), nil
}

// SetVoucher persists voucher for addr and emits the module's set-voucher event.
// It replaces any existing voucher for (addr, denom): voucher is the full new balance.
// Use when the caller has the exact intended amount (e.g. genesis import). To add to an
// existing balance without overwriting, use AddVoucher.
//
// Store value is the marshaled math.Int amount only (denom is in the key), matching bank-module
// balance storage semantics.
func (a *VouchersAssistant) SetVoucher(ctx sdk.Context, addr sdk.AccAddress, voucher sdk.Coin) (err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "SetVoucher")(&err)

	bz, err := voucher.Amount.Marshal()
	if err != nil {
		return err
	}

	a.keeper.GetVouchersStore(ctx).Set(voucherKey(voucher.Denom, addr), bz)
	a.keeper.EmitSetVoucherEvent(ctx, addr.String(), voucher)
	return nil
}

// DeleteVoucher removes the voucher for addr / denom and emits the module's delete-voucher event.
// It is a no-op if no voucher exists.
func (a *VouchersAssistant) DeleteVoucher(ctx sdk.Context, addr sdk.AccAddress, denom string) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "DeleteVoucher")()

	key := voucherKey(denom, addr)
	store := a.keeper.GetVouchersStore(ctx)
	if store.Has(key) {
		store.Delete(key)
		a.keeper.EmitDeleteVoucherEvent(ctx, addr.String(), denom)
	}
}

// AddVoucher adds coin to any existing voucher for (addr, denom) and persists the result.
// When no voucher exists for that pair this behaves like SetVoucher with that coin alone.
// Use this when crediting an incremental amount (e.g. repeated failed sends); using SetVoucher
// in those cases would overwrite prior claims.
func (a *VouchersAssistant) AddVoucher(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) (err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "AddVoucher")(&err)

	existing, err := a.GetVoucher(ctx, coin.Denom, addr)
	if err != nil {
		return err
	}
	return a.SetVoucher(ctx, addr, existing.Add(coin))
}

// ClaimVoucher sends the outstanding voucher amount from the module account to receiver and
// then deletes the voucher. Returns ErrVoucherNotFound if no voucher exists.
// If the bank transfer fails the voucher is left untouched so the user can retry.
func (a *VouchersAssistant) ClaimVoucher(ctx sdk.Context, receiver sdk.AccAddress, denom string) (err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "ClaimVoucher")(&err)

	voucher, err := a.GetVoucher(ctx, denom, receiver)
	if err != nil {
		return err
	}
	if voucher.IsZero() {
		return ErrVoucherNotFound
	}

	if err := a.bankKeeper.SendCoinsFromModuleToAccount(ctx, a.keeper.ModuleName(), receiver, sdk.NewCoins(voucher)); err != nil {
		return err
	}

	a.DeleteVoucher(ctx, receiver, denom)
	return nil
}

// GetAllVouchers returns every outstanding voucher across all denoms and addresses.
func (a *VouchersAssistant) GetAllVouchers(ctx sdk.Context) (result []vouchertypes.AddressVoucher, err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "GetAllVouchers")(&err)

	const addrLen = 20
	keySuffixLen := addrLen + len(delim)
	chaintypes.IterateSafe(a.keeper.GetVouchersStore(ctx).Iterator(nil, nil), func(k, v []byte) bool {
		if len(k) <= keySuffixLen {
			err = errors.Wrap(sdkerrors.ErrInvalidRequest, "voucher key too short")
			return true
		}
		denom := string(k[:len(k)-keySuffixLen])
		address := sdk.AccAddress(k[len(k)-addrLen:])
		var amount math.Int
		if err = amount.Unmarshal(v); err != nil {
			return true
		}
		result = append(result, vouchertypes.AddressVoucher{
			Address: address.String(),
			Voucher: sdk.NewCoin(denom, amount),
		})
		return false
	})
	return result, err
}

// GetVouchersForDenom returns every outstanding voucher for the given denom.
func (a *VouchersAssistant) GetVouchersForDenom(ctx sdk.Context, denom string) (result []vouchertypes.AddressVoucher, err error) {
	defer a.keeper.Meter(ctx).FuncTiming(&ctx, "GetVouchersForDenom")(&err)

	denomStore := prefix.NewStore(a.keeper.GetVouchersStore(ctx), denomWithDelim(denom))
	chaintypes.IterateSafe(denomStore.Iterator(nil, nil), func(k, v []byte) bool {
		var amount math.Int
		if err = amount.Unmarshal(v); err != nil {
			return true
		}
		result = append(result, vouchertypes.AddressVoucher{
			Address: sdk.AccAddress(k).String(),
			Voucher: sdk.NewCoin(denom, amount),
		})
		return false
	})
	return result, err
}

// voucherKey constructs the store key: denom + "|" + address_bytes.
func voucherKey(denom string, addr sdk.AccAddress) []byte {
	return append(denomWithDelim(denom), addr.Bytes()...)
}

func denomWithDelim(denom string) []byte {
	return append([]byte(denom), delim...)
}
