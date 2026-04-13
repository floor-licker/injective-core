package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "insurance"
	StoreKey   = ModuleName
)

var (
	// Key for insurance prefixes
	InsuranceFundPrefixKey = []byte{0x02}

	// Key for insurance redemption prefixes
	RedemptionSchedulePrefixKey = []byte{0x03}

	GlobalShareDenomIdPrefixKey         = []byte{0x04, 0x00}
	GlobalRedemptionScheduleIdPrefixKey = []byte{0x05, 0x00}

	// Key for redemption schedule secondary index by (redeemer address, marketID)
	RedemptionScheduleByAddrPrefixKey = []byte{0x06}

	// Keys for failed redemption schedule prefixes
	FailedRedemptionSchedulePrefixKey         = []byte{0x07}
	GlobalFailedRedemptionScheduleIdPrefixKey = []byte{0x08, 0x00}

	VouchersKey = []byte{0x09}

	ParamsKey = []byte{0x10}
)

// GetRedemptionScheduleKey provides the key to store a single pending redemption
func GetRedemptionScheduleKey(redemptionID uint64, claimTime time.Time) []byte {
	key := RedemptionSchedulePrefixKey
	key = append(key, sdk.FormatTimeBytes(claimTime)...)
	key = append(key, sdk.Uint64ToBigEndian(redemptionID)...)
	return key
}

// GetRedemptionScheduleKey provides the key to store a single pending redemption
func (sh RedemptionSchedule) GetRedemptionScheduleKey() []byte {
	return GetRedemptionScheduleKey(sh.Id, sh.ClaimableRedemptionTime)
}

// GetRedemptionScheduleByAddrKey returns the secondary index key for looking up
// redemption schedules by redeemer address and market ID.
// Key format: [prefix][addr_len][addr_bytes][marketID_hex_bytes][claimTime][redemptionID]
func GetRedemptionScheduleByAddrKey(redeemer sdk.AccAddress, marketID string, redemptionID uint64, claimTime time.Time) []byte {
	key := GetRedemptionScheduleByAddrPrefix(redeemer, marketID)
	key = append(key, sdk.FormatTimeBytes(claimTime)...)
	key = append(key, sdk.Uint64ToBigEndian(redemptionID)...)
	return key
}

// GetRedemptionScheduleByAddrPrefix returns the prefix for iterating all
// redemption schedules for a specific redeemer + market.
// Key format: [prefix][addr_len][addr_bytes][marketID_hex_bytes]
func GetRedemptionScheduleByAddrPrefix(redeemer sdk.AccAddress, marketID string) []byte {
	key := GetRedemptionScheduleByAddrOnlyPrefix(redeemer)
	key = append(key, []byte(marketID)...)
	return key
}

// GetRedemptionScheduleByAddrOnlyPrefix returns the prefix for iterating all
// redemption schedules for a specific redeemer across all markets.
// Key format: [prefix][addr_len][addr_bytes]
func GetRedemptionScheduleByAddrOnlyPrefix(redeemer sdk.AccAddress) []byte {
	key := RedemptionScheduleByAddrPrefixKey
	key = append(key, byte(len(redeemer)))
	key = append(key, redeemer...)
	return key
}

// GetRedemptionScheduleByAddrKeyFromSchedule is a convenience method on RedemptionSchedule.
func (sh RedemptionSchedule) GetRedemptionScheduleByAddrKey() ([]byte, error) {
	redeemer, err := sdk.AccAddressFromBech32(sh.Redeemer)
	if err != nil {
		return nil, err
	}
	return GetRedemptionScheduleByAddrKey(redeemer, sh.MarketId, sh.Id, sh.ClaimableRedemptionTime), nil
}

// GetFailedRedemptionScheduleKey provides the key to store a single failed redemption
func GetFailedRedemptionScheduleKey(failedRedemptionID uint64) []byte {
	key := FailedRedemptionSchedulePrefixKey
	key = append(key, sdk.Uint64ToBigEndian(failedRedemptionID)...)
	return key
}

// GetFailedRedemptionScheduleKey provides the key for a FailedRedemptionSchedule instance
func (f FailedRedemptionSchedule) GetFailedRedemptionScheduleKey() []byte {
	return GetFailedRedemptionScheduleKey(f.Id)
}
