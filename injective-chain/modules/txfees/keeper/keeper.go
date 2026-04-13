package keeper

import (
	"context"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/InjectiveLabs/metrics/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	mempool1559 "github.com/InjectiveLabs/injective-core/injective-chain/modules/txfees/keeper/mempool-1559"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/txfees/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	consensusKeeper types.ConsensusKeeper
	dataDir         string
	CurFeeState     *mempool1559.FeeState

	authority string

	meter metrics.Meter
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	consensusKeeper types.ConsensusKeeper,
	dataDir string,
	authority string,
) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		consensusKeeper: consensusKeeper,
		dataDir:         dataDir,
		// Initialize the EIP state with the default values. They will be updated in the BeginBlocker.
		CurFeeState: mempool1559.DefaultFeeState(),
		authority:   authority,
	}
}

func (*Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

func (k *Keeper) Meter(ctx context.Context) metrics.Meter {
	if k.meter == nil {
		k.meter = sdk.UnwrapSDKContext(ctx).Meter().SubMeter(types.ModuleName, metrics.Tag("svc", types.ModuleName))
	}

	return k.meter
}

// GetConsParams returns the current consensus parameters from the consensus params store.
func (k *Keeper) GetConsParams(ctx sdk.Context) (*consensustypes.QueryParamsResponse, error) {
	return k.consensusKeeper.Params(ctx, &consensustypes.QueryParamsRequest{})
}
