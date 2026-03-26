package v1dot18dot3beta

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/app/upgrades"
)

const (
	UpgradeVersion = "v1.18.3-beta"
)

func StoreUpgrades() storetypes.StoreUpgrades {
	return storetypes.StoreUpgrades{
		Added:   nil,
		Renamed: nil,
		Deleted: nil,
	}
}

func UpgradeSteps() []*upgrades.UpgradeHandlerStep {
	return []*upgrades.UpgradeHandlerStep{}
}
