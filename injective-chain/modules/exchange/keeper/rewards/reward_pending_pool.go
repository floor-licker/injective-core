package rewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
)

// GetAllCampaignRewardPendingPools gets all campaign reward pools
func (k TradingKeeper) GetAllCampaignRewardPendingPools(ctx sdk.Context) []*v2.CampaignRewardPool {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetAllCampaignRewardPendingPools")()

	rewardPools := make([]*v2.CampaignRewardPool, 0)
	k.IterateCampaignRewardPendingPools(ctx, false, func(pool *v2.CampaignRewardPool) (stop bool) {
		rewardPools = append(rewardPools, pool)
		return false
	})

	return rewardPools
}

// GetFirstCampaignRewardPendingPool gets the first campaign reward pool.
func (k TradingKeeper) GetFirstCampaignRewardPendingPool(ctx sdk.Context) (rewardPool *v2.CampaignRewardPool) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetFirstCampaignRewardPendingPool")()

	appendPool := func(pool *v2.CampaignRewardPool) (stop bool) {
		rewardPool = pool
		return true
	}

	k.IterateCampaignRewardPendingPools(ctx, false, appendPool)
	return rewardPool
}
