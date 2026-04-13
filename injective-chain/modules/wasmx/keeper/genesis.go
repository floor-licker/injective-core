package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "InitGenesis")(&err)

	k.SetParams(ctx, data.Params)
	for _, contract := range data.RegisteredContracts {
		address, err := sdk.AccAddressFromBech32(contract.Address)
		if err != nil {
			err = fmt.Errorf("error in contract address: %s: %w", contract.Address, err)
			panic(err)
		}
		k.SetContract(ctx, address, *contract.RegisteredContract)
	}

	k.CreateModuleAccount(ctx)
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	defer k.Meter(ctx).FuncTiming(&ctx, "ExportGenesis")()

	return &types.GenesisState{
		Params:              k.GetParams(ctx),
		RegisteredContracts: k.GetAllRegisteredContracts(ctx),
	}
}
