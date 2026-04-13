package auction

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	auctionkeeper "github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/keeper"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/auction/types"
)

func InitGenesis(ctx sdk.Context, keeper *auctionkeeper.Keeper, data types.GenesisState) {
	var err error
	defer keeper.Meter(ctx).FuncTiming(&ctx, "InitGenesis")(&err)

	keeper.SetParams(ctx, data.Params)

	// load highest bidder
	keeper.DeleteBid(ctx)
	if data.HighestBid != nil {
		keeper.SetBid(ctx, data.HighestBid.Bidder, data.HighestBid.Amount)
	}

	// load auction round
	keeper.SetAuctionRound(ctx, data.AuctionRound)

	// set ending time stamp for this round
	if data.AuctionEndingTimestamp == 0 {
		keeper.InitEndingTimeStamp(ctx)
	} else {
		keeper.SetEndingTimeStamp(ctx, data.AuctionEndingTimestamp)
	}

	if data.LastAuctionResult != nil {
		keeper.SetLastAuctionResult(ctx, *data.LastAuctionResult)
	}

	// Genesis carries exact voucher balances; SetVoucher restores them (AddVoucher would be wrong here).
	for _, av := range data.Vouchers {
		addr, err := sdk.AccAddressFromBech32(av.Address)
		if err != nil {
			err = fmt.Errorf("invalid voucher address in genesis: %w", err)
			panic(err)
		}
		err = keeper.SetVoucher(ctx, addr, av.Voucher)
		if err != nil {
			panic(err)
		}
	}

	keeper.CreateModuleAccount(ctx)
}

func ExportGenesis(ctx sdk.Context, k *auctionkeeper.Keeper) *types.GenesisState {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "ExportGenesis")(&err)

	avs, err := k.GetAllVouchers(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:                 k.GetParams(ctx),
		AuctionRound:           k.GetAuctionRound(ctx),
		HighestBid:             k.GetHighestBid(ctx),
		AuctionEndingTimestamp: k.GetEndingTimeStamp(ctx),
		LastAuctionResult:      k.GetLastAuctionResult(ctx),
		Vouchers:               avs,
	}
}
