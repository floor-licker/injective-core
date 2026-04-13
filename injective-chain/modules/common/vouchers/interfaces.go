package vouchers

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	metrics "github.com/InjectiveLabs/metrics/v2"
)

// VoucherKeeper is the interface a module's keeper must implement to work with VouchersAssistant.
// It provides the store access, the module name for the module account, event emission hooks,
// and the module's metrics meter so the assistant can instrument its operations.
type VoucherKeeper interface {
	GetVouchersStore(ctx sdk.Context) storetypes.KVStore
	ModuleName() string
	EmitSetVoucherEvent(ctx sdk.Context, addr string, voucher sdk.Coin)
	EmitDeleteVoucherEvent(ctx sdk.Context, addr string, denom string)
	Meter(ctx context.Context) metrics.Meter
}

// BankKeeper defines the bank operations required by VouchersAssistant to release claimed vouchers.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
