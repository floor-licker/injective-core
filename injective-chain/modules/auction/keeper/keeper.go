package keeper

import (
	"context"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/InjectiveLabs/metrics/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

// Keeper of this module maintains collections of auction.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper authkeeper.AccountKeeper
	bankKeeper    types.BankKeeper

	meter metrics.Meter

	authority string

	vouchersAssistant *vouchers.VouchersAssistant
}

// NewKeeper creates new instances of the auction Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	authority string,
) *Keeper {
	k := &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
		authority:     authority,
	}
	k.vouchersAssistant = vouchers.NewVouchersAssistant(k, bk)
	return k
}

// GetVouchersStore returns the KV store prefixed for voucher storage (satisfies vouchers.VoucherKeeper).
func (k Keeper) GetVouchersStore(ctx sdk.Context) storetypes.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.VouchersKey)
}

// ModuleName returns the auction module name (satisfies vouchers.VoucherKeeper).
func (Keeper) ModuleName() string { return types.ModuleName }

// EmitSetVoucherEvent emits EventSetVoucher (satisfies vouchers.VoucherKeeper).
func (Keeper) EmitSetVoucherEvent(ctx sdk.Context, addr string, voucher sdk.Coin) {
	if err := ctx.EventManager().EmitTypedEvent(&types.EventSetVoucher{
		Addr:    addr,
		Voucher: voucher,
	}); err != nil {
		ctx.Logger().Error("failed to emit EventSetVoucher", "addr", addr, "voucher", voucher, "err", err)
	}
}

// EmitDeleteVoucherEvent emits EventSetVoucher with a zero coin to signal deletion (satisfies vouchers.VoucherKeeper).
func (Keeper) EmitDeleteVoucherEvent(ctx sdk.Context, addr, denom string) {
	if err := ctx.EventManager().EmitTypedEvent(&types.EventSetVoucher{
		Addr:    addr,
		Voucher: types.NewEmptyVoucher(denom),
	}); err != nil {
		ctx.Logger().Error("failed to emit EventSetVoucher (delete)", "addr", addr, "denom", denom, "err", err)
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

func (k *Keeper) Meter(ctx context.Context) metrics.Meter {
	if k.meter == nil {
		k.meter = sdk.UnwrapSDKContext(ctx).Meter().SubMeter(types.ModuleName, metrics.Tag("svc", types.ModuleName))
	}

	return k.meter
}

func (k *Keeper) GetStore(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.storeKey)
}

// CreateModuleAccount creates the main module account and the auction fees subaccount (for ante fees).
// The fees subaccount is a derived module address that receives tx fees so they do not affect the ongoing round.
func (k *Keeper) CreateModuleAccount(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "CreateModuleAccount")()

	baseAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Burner)
	moduleAcc := (k.accountKeeper.NewAccount(ctx, baseAcc)).(sdk.ModuleAccountI) // set the account number
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)

	k.EnsureAuctionFeesSubaccount(ctx)
}

// EnsureAuctionFeesSubaccount creates the auction module's fee collector subaccount if it does not exist.
func (k *Keeper) EnsureAuctionFeesSubaccount(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "EnsureAuctionFeesSubaccount")()

	if !k.accountKeeper.HasAccount(ctx, types.AuctionFeesSubaccountAddress) {
		baseAccount := authtypes.NewBaseAccountWithAddress(types.AuctionFeesSubaccountAddress)
		feesSubacc := k.accountKeeper.NewAccount(ctx, baseAccount)
		k.accountKeeper.SetAccount(ctx, feesSubacc)
	}
}

// SweepFeesSubaccountToModule transfers INJ from the auction fees subaccount (ante fees)
// to the auction module account. Only the native INJ denom is swept; any other denom
// in the subaccount is ignored. We only expect INJ in the fees subaccount; refusing to
// sweep other denoms reduces attack surface (e.g. unexpected or malicious tokens with restricted permissions).
func (k *Keeper) SweepFeesSubaccountToModule(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SweepFeesSubaccountToModule")()

	injBalance := k.bankKeeper.GetBalance(ctx, types.AuctionFeesSubaccountAddress, chaintypes.InjectiveCoin)
	if injBalance.Amount.IsZero() {
		return
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, types.AuctionFeesSubaccountAddress, types.ModuleName, sdk.NewCoins(injBalance)); err != nil {
		// It is not expected that transferring the native coin (INJ) would produce any error really.
		// But if it does, launching the new auction round can't fail. We prefer to continue the new round initialization
		// without the accumulated ante fees.
		k.Logger(ctx).Error(
			"SweepFeesSubaccountToModule failed to transfer INJ to auction module",
			"amount", injBalance.Amount.String(),
			"error", err,
		)
	}
}
