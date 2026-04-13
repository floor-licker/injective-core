package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func (gs GenesisState) Validate() error {
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
		Params:                 DefaultParams(),
		AuctionRound:           0,
		HighestBid:             nil,
		AuctionEndingTimestamp: 0,
	}
}
