package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func (gs GenesisState) Validate() error {
	if gs.NextRedemptionScheduleId == 0 {
		return errors.New("NextRedemptionScheduleId should NOT be zero")
	}
	if gs.NextShareDenomId == 0 {
		return errors.New("NextShareDenomId should NOT be zero")
	}
	// NextFailedRedemptionScheduleId may be 0 in legacy genesis (field omitted); treated as 1 in InitGenesis

	for _, s := range gs.RedemptionSchedule {
		if s.Id >= gs.NextRedemptionScheduleId {
			return errors.New("RedemptionSchedule id must be less than NextRedemptionScheduleId")
		}
	}

	nextFailedID := gs.NextFailedRedemptionScheduleId
	if nextFailedID == 0 {
		nextFailedID = 1
	}
	for _, f := range gs.FailedRedemptionSchedules {
		if f.Id >= nextFailedID {
			return errors.New("FailedRedemptionSchedule id must be less than NextFailedRedemptionScheduleId")
		}
	}

	if err := gs.Params.Validate(); err != nil {
		return err
	}
	for _, av := range gs.Vouchers {
		if _, err := sdk.AccAddressFromBech32(av.Address); err != nil {
			return fmt.Errorf("invalid voucher address %q: %w", av.Address, err)
		}
		if err := av.Voucher.Validate(); err != nil {
			return fmt.Errorf("invalid voucher coin for address %q: %w", av.Address, err)
		}
	}
	return nil
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:                         DefaultParams(),
		NextShareDenomId:               1,
		NextRedemptionScheduleId:       1,
		NextFailedRedemptionScheduleId: 1,
		RedemptionSchedule:             []RedemptionSchedule{},
		InsuranceFunds:                 []InsuranceFund{},
		FailedRedemptionSchedules:      []FailedRedemptionSchedule{},
	}
}
