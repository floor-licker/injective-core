package rewards

import (
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k TradingKeeper) SetTradingRewardsMarketQualificationForAllQualifyingMarkets(
	ctx sdk.Context,
	campaignInfo *v2.TradingRewardCampaignInfo,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "SetTradingRewardsMarketQualificationForAllQualifyingMarkets")()

	marketIDQuoteDenoms := k.GetAllMarketIDsWithQuoteDenoms(ctx)

	quoteDenomMap := make(map[string]struct{})
	for _, quoteDenom := range campaignInfo.QuoteDenoms {
		quoteDenomMap[quoteDenom] = struct{}{}
	}

	for _, m := range marketIDQuoteDenoms {
		if _, ok := quoteDenomMap[m.QuoteDenom]; ok {
			k.SetTradingRewardsMarketQualification(ctx, m.MarketID, true)
		}
	}

	for _, marketID := range campaignInfo.DisqualifiedMarketIds {
		k.SetTradingRewardsMarketQualification(ctx, common.HexToHash(marketID), false)
	}
}

// DeleteAllTradingRewardsMarketQualifications deletes the trading reward qualifications for all markets
func (k TradingKeeper) DeleteAllTradingRewardsMarketQualifications(ctx sdk.Context) {
	defer k.Meter(ctx).FuncTiming(&ctx, "DeleteAllTradingRewardsMarketQualifications")()

	marketIDs, _ := k.GetAllTradingRewardsMarketQualification(ctx)
	for _, marketID := range marketIDs {
		k.DeleteTradingRewardsMarketQualification(ctx, marketID)
	}
}

// GetAllTradingRewardsMarketQualification gets all market qualification statuses
func (k TradingKeeper) GetAllTradingRewardsMarketQualification(ctx sdk.Context) ([]common.Hash, []bool) {
	defer k.Meter(ctx).FuncTiming(&ctx, "GetAllTradingRewardsMarketQualification")()

	marketIDs := make([]common.Hash, 0)
	isQualified := make([]bool, 0)

	appendQualification := func(m common.Hash, q bool) (stop bool) {
		marketIDs = append(marketIDs, m)
		isQualified = append(isQualified, q)
		return false
	}

	k.IterateTradingRewardsMarketQualifications(ctx, appendQualification)

	return marketIDs, isQualified
}

func (k TradingKeeper) CheckQuoteAndSetTradingRewardQualification(
	ctx sdk.Context,
	marketID common.Hash,
	quoteDenom string,
) {
	defer k.Meter(ctx).FuncTiming(&ctx, "CheckQuoteAndSetTradingRewardQualification")()

	if campaign := k.GetCampaignInfo(ctx); campaign != nil {
		disqualified := false
		for _, disqualifiedMarketID := range campaign.DisqualifiedMarketIds {
			if marketID == common.HexToHash(disqualifiedMarketID) {
				disqualified = true
			}
		}

		if disqualified {
			k.SetTradingRewardsMarketQualification(ctx, marketID, false)
			return
		}

		for _, q := range campaign.QuoteDenoms {
			if quoteDenom == q {
				k.SetTradingRewardsMarketQualification(ctx, marketID, true)
				break
			}
		}
	}
}
