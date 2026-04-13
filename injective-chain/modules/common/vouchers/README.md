# Vouchers Assistant

`VouchersAssistant` is a reusable component that centralises all voucher CRUD logic. Any Injective module can gain full voucher support by implementing a small interface and wiring in this component — no business logic needs to be duplicated.

## Concept

A voucher represents a token amount that could not be delivered to a recipient at the time of a transfer. Instead of failing the operation, the sending module stores a voucher on behalf of the recipient. The recipient can claim the voucher at any time in the future.

Typical trigger: a bank transfer inside consensus-critical code fails (e.g. due to a send restriction). The module intercepts the failure, accumulates the amount into a voucher, and redirects the bank transfer to its own module account so the block does not fail. The voucher is stored keyed by `(denom, recipient address)`.

## Storage layout

`VouchersAssistant` stores every voucher inside the KV store slice that the module provides via `GetVouchersStore`. The key layout inside that store is:

```
denom + "|" + address_bytes
```

Because the assistant only writes within the store returned by `GetVouchersStore`, the module's own prefix is applied exactly once by that method (e.g. via `prefix.NewStore`). There is no double-prefixing.

The value is the binary encoding of `math.Int` for the amount only (same idea as bank-module balances: denom lives in the key, not duplicated in the value).

## API

```go
// NewVouchersAssistant wires the assistant to a module keeper and bank keeper.
func NewVouchersAssistant(keeper VoucherKeeper, bankKeeper BankKeeper) *VouchersAssistant

// GetVoucher returns the outstanding coin for (denom, addr).
// Returns a zero-amount coin when no voucher exists.
func (a *VouchersAssistant) GetVoucher(ctx sdk.Context, denom string, addr sdk.AccAddress) (sdk.Coin, error)

// SetVoucher persists a voucher and emits the module's set-voucher event. It replaces any
// existing balance for (addr, denom); use for genesis or when the amount is fully known.
func (a *VouchersAssistant) SetVoucher(ctx sdk.Context, addr sdk.AccAddress, voucher sdk.Coin) error

// AddVoucher adds to the existing voucher for (addr, denom), or creates one if missing.
// Use when crediting incremental amounts (e.g. repeated failed sends); do not use SetVoucher
// for that or prior claims would be overwritten.
func (a *VouchersAssistant) AddVoucher(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error

// DeleteVoucher removes the voucher for (addr, denom) and emits the delete-voucher event.
// No-op if no voucher exists.
func (a *VouchersAssistant) DeleteVoucher(ctx sdk.Context, addr sdk.AccAddress, denom string)

// ClaimVoucher sends the outstanding amount from the module account to receiver and deletes
// the voucher. Returns ErrVoucherNotFound if no voucher exists. The voucher is left intact
// if the bank transfer fails so the receiver can retry.
func (a *VouchersAssistant) ClaimVoucher(ctx sdk.Context, receiver sdk.AccAddress, denom string) error

// GetAllVouchers returns every outstanding voucher across all denoms and addresses.
func (a *VouchersAssistant) GetAllVouchers(ctx sdk.Context) ([]AddressVoucher, error)

// GetVouchersForDenom returns every outstanding voucher for a specific denom.
func (a *VouchersAssistant) GetVouchersForDenom(ctx sdk.Context, denom string) ([]AddressVoucher, error)
```

### Errors

```go
// ErrVoucherNotFound is returned by ClaimVoucher when no voucher exists for
// the requested (denom, address) pair.
var ErrVoucherNotFound = errors.Register("vouchers", 1, "voucher not found")
```

## Adding voucher support to a new module

Follow the steps below. The `permissions` module (`injective-chain/modules/permissions/keeper/`) is the canonical reference implementation.

### Step 1 — Reserve a store prefix

In your module's `keeper/keys.go`, reserve a unique byte prefix for the voucher store:

```go
var vouchersKey = []byte{0xNN} // pick the next available prefix byte
```

### Step 2 — Implement `VoucherKeeper` on your keeper

Your keeper must satisfy the `vouchers.VoucherKeeper` interface:

```go
type VoucherKeeper interface {
    GetVouchersStore(ctx sdk.Context) storetypes.KVStore
    ModuleName() string
    EmitSetVoucherEvent(ctx sdk.Context, addr string, voucher sdk.Coin)
    EmitDeleteVoucherEvent(ctx sdk.Context, addr string, denom string)
    Meter(ctx context.Context) metrics.Meter
}
```

Add these methods to `keeper/keeper.go`:

```go
func (k Keeper) GetVouchersStore(ctx sdk.Context) storetypes.KVStore {
    return prefix.NewStore(ctx.KVStore(k.storeKey), vouchersKey)
}

func (k Keeper) ModuleName() string { return types.ModuleName }

func (k Keeper) EmitSetVoucherEvent(ctx sdk.Context, addr string, voucher sdk.Coin) {
    // nolint:errcheck //ignored on purpose
    ctx.EventManager().EmitTypedEvent(&types.EventSetVoucher{
        Addr:    addr,
        Voucher: voucher,
    })
}

func (k Keeper) EmitDeleteVoucherEvent(ctx sdk.Context, addr string, denom string) {
    // nolint:errcheck //ignored on purpose
    ctx.EventManager().EmitTypedEvent(&types.EventSetVoucher{
        Addr:    addr,
        Voucher: types.NewEmptyVoucher(denom),
    })
}
```

`Meter` is typically already implemented on the keeper. If not, add it:

```go
func (k *Keeper) Meter(ctx context.Context) metrics.Meter {
    if k.meter == nil {
        k.meter = sdk.UnwrapSDKContext(ctx).Meter().SubMeter(types.ModuleName, metrics.Tag("svc", types.ModuleName))
    }
    return k.meter
}
```

> The event types (`EventSetVoucher`, `NewEmptyVoucher`) are module-specific. Define them in your module's `types/` package.

### Step 3 — Add the assistant field and initialise it

**Important:** `NewKeeper` must return `*Keeper` (pointer), not `Keeper` (value). This is required because `VouchersAssistant` holds a `VoucherKeeper` reference back to the keeper. If the keeper were a value type, the pointer passed during construction would become invalid after `NewKeeper` returns.

```go
type Keeper struct {
    // ... existing fields ...
    vouchersAssistant *vouchers.VouchersAssistant
}

func NewKeeper(storeKey storetypes.StoreKey, bankKeeper types.BankKeeper, /* ... */) *Keeper {
    k := &Keeper{
        storeKey:   storeKey,
        bankKeeper: bankKeeper,
        // ... other fields ...
    }
    k.vouchersAssistant = vouchers.NewVouchersAssistant(k, bankKeeper)
    return k
}
```

`NewVouchersAssistant(k, bankKeeper)` receives `k` (the pointer just created), so the back-reference is always valid.

### Step 4 — Define proto messages

Each module must define its own proto messages for `MsgClaimVoucher` and the voucher queries. The assistant's business logic is shared, but gRPC routing is per-type-URL, so these cannot be reused across modules.

Minimum required protos:

```protobuf
// tx.proto
message MsgClaimVoucher {
    string sender = 1;
    string denom  = 2;
}
message MsgClaimVoucherResponse {}

// query.proto
message QueryVouchersRequest  { string denom = 1; }  // empty denom → all
message QueryVouchersResponse { repeated AddressVoucher vouchers = 1; }

message QueryVoucherRequest   { string denom = 1; string address = 2; }
message QueryVoucherResponse  { cosmos.base.v1beta1.Coin voucher = 1; }
```

Regenerate bindings with `make proto-gen` after editing `.proto` files.

### Step 5 — Delegate in the message server

```go
func (k msgServer) ClaimVoucher(c context.Context, msg *types.MsgClaimVoucher) (*types.MsgClaimVoucherResponse, error) {
    ctx := sdk.UnwrapSDKContext(c)
    defer k.Meter(ctx).FuncTiming(&ctx, "ClaimVoucher")()

    receiver := sdk.MustAccAddressFromBech32(msg.Sender)

    if err := k.vouchersAssistant.ClaimVoucher(ctx, receiver, msg.Denom); err != nil {
        return nil, err
    }

    return &types.MsgClaimVoucherResponse{}, nil
}
```

Do not wrap or replace the `vouchers.ErrVoucherNotFound` error — return it as-is so callers can inspect the exact error type.

### Step 6 — Delegate in the query server

```go
func (q queryServer) Vouchers(c context.Context, req *types.QueryVouchersRequest) (*types.QueryVouchersResponse, error) {
    ctx := sdk.UnwrapSDKContext(c)
    defer q.Meter(ctx).FuncTiming(&ctx, "Vouchers")()

    var (
        raw []voucherspkg.AddressVoucher
        err error
    )
    if req.Denom == "" {
        raw, err = q.vouchersAssistant.GetAllVouchers(ctx)
    } else {
        raw, err = q.vouchersAssistant.GetVouchersForDenom(ctx, req.Denom)
    }
    if err != nil {
        return nil, err
    }

    result := make([]*types.AddressVoucher, 0, len(raw))
    for _, v := range raw {
        v := v
        result = append(result, &types.AddressVoucher{Address: v.Address, Voucher: v.Voucher})
    }
    return &types.QueryVouchersResponse{Vouchers: result}, nil
}

func (q queryServer) Voucher(c context.Context, req *types.QueryVoucherRequest) (*types.QueryVoucherResponse, error) {
    ctx := sdk.UnwrapSDKContext(c)
    defer q.Meter(ctx).FuncTiming(&ctx, "Voucher")()

    addr, err := sdk.AccAddressFromBech32(req.Address)
    if err != nil {
        return nil, err
    }

    voucher, err := q.vouchersAssistant.GetVoucher(ctx, req.Denom, addr)
    if err != nil {
        return nil, err
    }

    return &types.QueryVoucherResponse{Voucher: voucher}, nil
}
```

The `voucherspkg` alias avoids a name collision between the common package and the module's own `types.AddressVoucher`:

```go
import voucherspkg "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers"
```

### Step 7 — Wire genesis import / export

```go
// InitGenesis
for _, v := range genState.Vouchers {
    address := sdk.MustAccAddressFromBech32(v.Address)
    if err := k.vouchersAssistant.SetVoucher(ctx, address, v.Voucher); err != nil {
        panic(err)
    }
}

// ExportGenesis
rawVouchers, err := k.vouchersAssistant.GetAllVouchers(ctx)
if err != nil {
    panic(err)
}
gs.Vouchers = make([]*types.AddressVoucher, 0, len(rawVouchers))
for _, v := range rawVouchers {
    v := v
    gs.Vouchers = append(gs.Vouchers, &types.AddressVoucher{Address: v.Address, Voucher: v.Voucher})
}
```

### Step 8 — Metrics

Every method in `VouchersAssistant` instruments itself via `a.keeper.Meter(ctx).FuncTiming(...)`. Any keeper method that wraps or extends assistant calls should add its own `FuncTiming` using its exact method name:

```go
func (k Keeper) GetVoucherForAddress(ctx sdk.Context, denom string, addr sdk.AccAddress) (sdk.Coin, error) {
    defer k.Meter(ctx).FuncTiming(&ctx, "GetVoucherForAddress")()
    return k.vouchersAssistant.GetVoucher(ctx, denom, addr)
}
```

This creates a natural timing hierarchy: the outer keeper method's span wraps the inner assistant method's span.

### Step 9 — Tests

Add unit tests for the assistant using the in-memory mock pattern from `assistant_test.go`:

- Implement a `mockKeeper` satisfying `VoucherKeeper` (store backed by an in-memory `CommitMultiStore`; `Meter` returns `metrics.NewNilMeter()`).
- Implement a `mockBankKeeper` that records sent transfers and can be configured to return an error.
- Call `vouchers.NewVouchersAssistant(mockKeeper, mockBankKeeper)` and test all methods.

For integration tests, verify the full claim flow: trigger a failed send → voucher is created → call `MsgClaimVoucher` → voucher is delivered and removed.

## Optional: fail-fast send interception

If your module uses a bank send restriction (registered via `bankKeeper.AppendSendRestriction`) and needs to reroute failed transfers to vouchers instead of aborting the block, add a helper like the permissions module's `rerouteToVoucherOnFail`:

```go
func (k Keeper) rerouteToVoucherOnFail(ctx sdk.Context, toAddr sdk.AccAddress, amount sdk.Coin, origErr error) (sdk.AccAddress, error) {
    defer k.Meter(ctx).FuncTiming(&ctx, "rerouteToVoucherOnFail")()

    // Only reroute inside consensus-critical paths (indicated by DoNotFailFastSendContextKey).
    if ctx.Value(baseapp.DoNotFailFastSendContextKey) == nil {
        return toAddr, origErr
    }

    if err := k.vouchersAssistant.AddVoucher(ctx, toAddr, amount); err != nil {
        return toAddr, errors.Wrapf(err, "can't add voucher for address, tried to reroute token send after error: %s", origErr.Error())
    }

    return authtypes.NewModuleAddress(k.ModuleName()), nil
}
```

This logic is intentionally kept in the module's own keeper (not in the assistant) because it is specific to the bank send-restriction mechanism.

## Reference implementation

| File | Role |
|---|---|
| `injective-chain/modules/permissions/keeper/keeper.go` | `VoucherKeeper` implementation + assistant wiring |
| `injective-chain/modules/permissions/keeper/keys.go` | `vouchersKey` store prefix |
| `injective-chain/modules/permissions/keeper/vouchers.go` | Module-level voucher helpers + `rerouteToVoucherOnFail` |
| `injective-chain/modules/permissions/keeper/msg_server.go` | `ClaimVoucher` delegation |
| `injective-chain/modules/permissions/keeper/grpc_query.go` | `Vouchers` / `Voucher` query delegation |
| `injective-chain/modules/permissions/keeper/genesis.go` | Genesis import / export |
| `injective-chain/modules/common/vouchers/assistant_test.go` | Unit test pattern for new modules |
