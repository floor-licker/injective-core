package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// ContractBlacklistListenerFunc adapts a function to ContractBlacklistListener.
type ContractBlacklistListenerFunc func(ctx sdk.Context, contract, user common.Address) error

// OnEnforcedRestrictionsEVMContractBlacklist implements ContractBlacklistListener.
func (f ContractBlacklistListenerFunc) OnEnforcedRestrictionsEVMContractBlacklist(ctx sdk.Context, contract, user common.Address) error {
	return f(ctx, contract, user)
}

// NewAccountContractBlacklistListener adapts a callback that only needs the blacklisted account.
func NewAccountContractBlacklistListener(fn func(ctx sdk.Context, user sdk.AccAddress) error) ContractBlacklistListener {
	return ContractBlacklistListenerFunc(func(ctx sdk.Context, _ common.Address, user common.Address) error {
		return fn(ctx, sdk.AccAddress(user.Bytes()))
	})
}
