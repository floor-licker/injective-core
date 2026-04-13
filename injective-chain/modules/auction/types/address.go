package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// AuctionFeesSubaccountAddress is the address of the auction module's fee collector subaccount (module derivation path).
// Ante fees are sent here so they do not affect the ongoing auction round; they are swept into the auction
// module at round end by the auction module (SweepFeesSubaccountToModule). Use this address for new integrations.
//
// Mainnet (inj): inj18kc70l78dh9fjdn2nkvg4ezy73dcmyrgg234m4729q3mwmy2nzes3a8gx9
var AuctionFeesSubaccountAddress = sdk.AccAddress(address.Module(ModuleName, []byte(AuctionFeesSubaccountDerivationKey)))
