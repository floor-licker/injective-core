package insurance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (am AppModule) EndBlocker(ctx sdk.Context) {
	defer am.keeper.Meter(ctx).FuncTiming(&ctx, "EndBlocker")()
	// call automatic withdraw keeper function
	am.keeper.WithdrawAllMaturedRedemptions(ctx) //nolint:errcheck // ok
}
