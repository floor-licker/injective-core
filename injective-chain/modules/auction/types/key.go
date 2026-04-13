package types

const (
	ModuleName = "auction"
	StoreKey   = ModuleName
	TStoreKey  = "transient_auction"
	// AuctionFeesSubaccountDerivationKey is the derivation key for the auction module's fee collector subaccount.
	AuctionFeesSubaccountDerivationKey = "fees"
)

var (
	// Keys for store prefixes
	BidsKey              = []byte{0x01}
	AuctionRoundKey      = []byte{0x03}
	KeyEndingTimeStamp   = []byte{0x04}
	KeyLastAuctionResult = []byte{0x05}
	VouchersKey          = []byte{0x06}
	ParamsKey            = []byte{0x10}
)
