package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/permissions/types"
)

// InitGenesis initializes the permissions module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "InitGenesis")(&err)

	err = genState.Params.Validate()
	if err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)

	for idx := range genState.Namespaces {
		err = k.createNamespace(ctx, genState.Namespaces[idx])
		if err != nil {
			panic(err)
		}
	}

	// Genesis carries exact voucher balances; SetVoucher restores them (AddVoucher would be wrong here).
	for _, v := range genState.Vouchers {
		address, err := sdk.AccAddressFromBech32(v.Address)
		if err != nil {
			panic(err)
		}
		err = k.vouchersAssistant.SetVoucher(ctx, address, v.Voucher)
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the permissions module's exported genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "ExportGenesis")(&err)

	namespaces, err := k.GetAllNamespaces(ctx)
	if err != nil {
		panic(err)
	}

	gs := &types.GenesisState{Params: k.GetParams(ctx)}
	for _, ns := range namespaces {
		gs.Namespaces = append(gs.Namespaces, *ns)
	}

	rawVouchers, err := k.vouchersAssistant.GetAllVouchers(ctx)
	if err != nil {
		panic(err)
	}

	gs.Vouchers = rawVouchers

	return gs
}
