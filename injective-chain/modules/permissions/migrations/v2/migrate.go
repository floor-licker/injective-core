package v2

import (
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var vouchersKey = []byte{0x09}

// Migrate converts permissions voucher values from legacy proto-encoded sdk.Coin
// (as written by master) to amount-only math.Int bytes (VouchersAssistant / bank-style).
func Migrate(store storetypes.KVStore) error {
	vouchersStore := prefix.NewStore(store, vouchersKey)

	type entry struct {
		key, value []byte
	}
	var entries []entry
	iter := vouchersStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		entries = append(entries, entry{
			key:   append([]byte(nil), iter.Key()...),
			value: append([]byte(nil), iter.Value()...),
		})
	}
	iter.Close()

	for _, e := range entries {
		var coin sdk.Coin
		if err := proto.Unmarshal(e.value, &coin); err != nil {
			var amt math.Int
			if err2 := amt.Unmarshal(e.value); err2 == nil {
				continue
			}
			return fmt.Errorf("failed to unmarshal legacy voucher at key %x: %w", e.key, err)
		}

		bz, err := coin.Amount.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal voucher amount: %w", err)
		}
		vouchersStore.Set(e.key, bz)
	}

	return nil
}
