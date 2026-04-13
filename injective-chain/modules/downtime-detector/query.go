package downtimedetector

import (
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/downtime-detector/types"
)

func (k *Keeper) RecoveredSinceDowntimeOfLength(ctx sdk.Context, downtime types.Downtime, recoveryDuration time.Duration) (recovered bool, err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "RecoveredSinceDowntimeOfLength")(&err)

	lastDowntime, err := k.GetLastDowntimeOfLength(ctx, downtime)
	if err != nil {
		return false, err
	}
	if recoveryDuration == time.Duration(0) {
		return false, errors.New("invalid recovery duration of 0")
	}
	// Check if current time < lastDowntime + recovery duration
	// if LTE, then we have not waited recovery duration.
	if ctx.BlockTime().Before(lastDowntime.Add(recoveryDuration)) {
		return false, nil
	}
	return true, nil
}
