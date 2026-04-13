package auction

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// NewAuctionProposalHandler returns a governance handler for auction proposal types.
// No auction-specific proposal types are currently defined; any routed content
// returns ErrUnknownRequest.
func NewAuctionProposalHandler() govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		default:
			return errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized auction proposal content type: %T", c)
		}
	}
}
