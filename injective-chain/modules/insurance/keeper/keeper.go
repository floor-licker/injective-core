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

	exchangekeeper "github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/keeper"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/insurance/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

// Keeper of this module maintains collections of insurance.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper  authkeeper.AccountKeeper
	bankKeeper     types.BankKeeper
	exchangeKeeper *exchangekeeper.Keeper

	meter metrics.Meter

	authority string

	vouchersAssistant *vouchers.VouchersAssistant
}

// NewKeeper creates new instances of the insurance Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	ek *exchangekeeper.Keeper,
	authority string,
) *Keeper {
	k := &Keeper{
		storeKey:       storeKey,
		cdc:            cdc,
		accountKeeper:  ak,
		bankKeeper:     bk,
		exchangeKeeper: ek,
		authority:      authority,
	}
	k.vouchersAssistant = vouchers.NewVouchersAssistant(k, bk)
	return k
}

// GetVouchersStore returns the KV store prefixed for voucher storage (satisfies vouchers.VoucherKeeper).
func (k Keeper) GetVouchersStore(ctx sdk.Context) storetypes.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.VouchersKey)
}

// ModuleName returns the insurance module name (satisfies vouchers.VoucherKeeper).
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

func (k *Keeper) GetStore(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.storeKey)
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

// CreateModuleAccount creates a module account with minting and burning capabilities
func (k *Keeper) CreateModuleAccount(ctx sdk.Context) {
	baseAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter, authtypes.Burner)
	moduleAcc := (k.accountKeeper.NewAccount(ctx, baseAcc)).(sdk.ModuleAccountI) // set the account number
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

func (k *Keeper) SetExchangeKeeper(ek *exchangekeeper.Keeper) {
	k.exchangeKeeper = ek
}

// BackfillRedemptionScheduleAddrIndex populates the (redeemer, marketID) secondary
// index for all existing redemption schedules. This must be run once as a state
// migration when upgrading to the version that introduced the secondary index.
// Errors on individual schedules are logged and skipped to avoid blocking the upgrade.
func (k *Keeper) BackfillRedemptionScheduleAddrIndex(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "BackfillRedemptionScheduleAddrIndex")()

	store := ctx.KVStore(k.storeKey)

	chaintypes.IterateSafe(k.globalRedemptionIterator(ctx), func(_, value []byte) bool {
		schedule := k.unmarshalRedemptionSchedule(value)
		if schedule == nil {
			k.Logger(ctx).Error("skipping redemption schedule: unmarshal failure during backfill")
			return false
		}

		addrKey, err := schedule.GetRedemptionScheduleByAddrKey()
		if err != nil {
			k.Logger(ctx).Error("skipping redemption schedule: failed to compute addr index key", "id", schedule.Id, "error", err)
			return false
		}
		store.Set(addrKey, []byte{})
		return false
	})
}
