package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

const (
	RouterKey = ModuleName

	TypeMsgBid          = "bid"
	TypeMsgUpdateParams = "updateParams"
	TypeMsgClaimVoucher = "claimVoucher"
)

var (
	_ sdk.Msg = &MsgBid{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgClaimVoucher{}
)

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	if err := msg.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgBid) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgBid) Type() string { return TypeMsgBid }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !msg.BidAmount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.BidAmount.String())
	}

	if !msg.BidAmount.IsPositive() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, msg.BidAmount.String())
	}

	if msg.BidAmount.Denom != chaintypes.InjectiveCoin {
		return errors.Wrap(ErrBidInvalid, msg.BidAmount.Denom)
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgBid) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface.
func (MsgClaimVoucher) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (MsgClaimVoucher) Type() string { return TypeMsgClaimVoucher }

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgClaimVoucher) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}
	return nil
}

// GetSignBytes implements the sdk.Msg interface.
func (msg *MsgClaimVoucher) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgClaimVoucher) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// NewEmptyVoucher returns a zero-amount coin for the given denom.
// Used by EmitDeleteVoucherEvent to signal voucher deletion via an EventSetVoucher event.
func NewEmptyVoucher(denom string) sdk.Coin {
	return sdk.NewInt64Coin(denom, 0)
}
