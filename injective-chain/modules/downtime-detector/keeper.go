package downtimedetector

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/downtime-detector/types"
	"github.com/InjectiveLabs/metrics/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	meter    metrics.Meter
}

func NewKeeper(storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{storeKey: storeKey}
}

func (k *Keeper) Meter(ctx context.Context) metrics.Meter {
	if k.meter == nil {
		k.meter = sdk.UnwrapSDKContext(ctx).Meter().SubMeter(types.ModuleName, metrics.Tag("svc", types.ModuleName))
	}

	return k.meter
}
