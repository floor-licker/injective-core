package insurance

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/insurance/keeper"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/insurance/types"
)

// InitGenesis init state of module
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, data types.GenesisState) {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "InitGenesis")(&err)

	k.SetParams(ctx, data.Params)
	for i := range data.InsuranceFunds {
		k.SetInsuranceFund(ctx, &data.InsuranceFunds[i])
	}
	for _, schedule := range data.RedemptionSchedule {
		k.SetRedemptionSchedule(ctx, schedule)
	}
	k.SetNextShareDenomId(ctx, data.NextShareDenomId)
	k.SetNextRedemptionScheduleId(ctx, data.NextRedemptionScheduleId)
	for _, failed := range data.FailedRedemptionSchedules {
		k.SetFailedRedemptionSchedule(ctx, failed)
	}
	nextFailedID := data.NextFailedRedemptionScheduleId
	if nextFailedID == 0 {
		nextFailedID = 1
	}
	k.SetNextFailedRedemptionScheduleId(ctx, nextFailedID)

	// Genesis carries exact voucher balances; SetVoucher restores them (AddVoucher would be wrong here).
	for _, av := range data.Vouchers {
		addr, err := sdk.AccAddressFromBech32(av.Address)
		if err != nil {
			err = fmt.Errorf("invalid voucher address in genesis: %w", err)
			panic(err)
		}
		err = k.SetVoucher(ctx, addr, av.Voucher)
		if err != nil {
			panic(err)
		}
	}

	k.CreateModuleAccount(ctx)
}

// ExportGenesis export the state of module
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	var err error
	defer k.Meter(ctx).FuncTiming(&ctx, "ExportGenesis")(&err)

	avs, err := k.GetAllVouchers(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:                         k.GetParams(ctx),
		InsuranceFunds:                 k.GetAllInsuranceFunds(ctx),
		RedemptionSchedule:             k.GetAllInsuranceFundRedemptions(ctx),
		NextShareDenomId:               k.ExportNextShareDenomId(ctx),
		NextRedemptionScheduleId:       k.ExportNextRedemptionScheduleId(ctx),
		FailedRedemptionSchedules:      k.GetAllFailedRedemptionSchedules(ctx),
		NextFailedRedemptionScheduleId: k.ExportNextFailedRedemptionScheduleId(ctx),
		Vouchers:                       avs,
	}
}
