package wasmx

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/keeper"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/types"
)

// NewWasmxProposalHandler creates a governance handler to manage new wasmx proposal types.
func NewWasmxProposalHandler(k keeper.Keeper, wasmProposalHandler govtypes.Handler) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ContractRegistrationRequestProposal:
			return handleContractRegistrationRequestProposal(ctx, k, c)
		case *types.BatchContractRegistrationRequestProposal:
			return handleBatchContractRegistrationRequestProposal(ctx, k, c)
		case *types.BatchContractDeregistrationProposal:
			return handleBatchContractDeregistrationProposal(ctx, k, c)
		case *types.BatchStoreCodeProposal:
			return handleBatchStoreCodeProposal(ctx, k, c, wasmProposalHandler)
		default:
			return errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasmx proposal content type: %T", c)
		}
	}
}

func handleContractRegistrationRequestProposal(ctx sdk.Context, k keeper.Keeper, p *types.ContractRegistrationRequestProposal) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleContractRegistrationRequestProposal")(&err)

	if err := p.ValidateBasic(); err != nil {
		return err
	}

	params := k.GetParams(ctx)
	return k.HandleContractRegistration(ctx, params, p.ContractRegistrationRequest)
}

func handleBatchContractRegistrationRequestProposal(ctx sdk.Context, k keeper.Keeper, p *types.BatchContractRegistrationRequestProposal) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleBatchContractRegistrationRequestProposal")(&err)

	if err := p.ValidateBasic(); err != nil {
		return err
	}

	params := k.GetParams(ctx)

	for _, req := range p.ContractRegistrationRequests {
		if err := k.HandleContractRegistration(ctx, params, req); err != nil {
			return err
		}
	}

	return nil
}

func handleBatchContractDeregistrationProposal(ctx sdk.Context, k keeper.Keeper, p *types.BatchContractDeregistrationProposal) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleBatchContractDeregistrationProposal")(&err)

	if err := p.ValidateBasic(); err != nil {
		return err
	}

	for _, contract := range p.Contracts {
		contractAddress, err := sdk.AccAddressFromBech32(contract)
		if err != nil {
			return err
		}

		if err := k.DeregisterContract(ctx, contractAddress); err != nil {
			if sdkerrors.ErrNotFound.Is(err) {
				continue // no need to break processing if contract is not registered, just skip
			} else {
				return err
			}
		}
	}

	return nil
}

func handleBatchStoreCodeProposal(ctx sdk.Context, k keeper.Keeper, p *types.BatchStoreCodeProposal, wasmProposalHandler govtypes.Handler) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "handleBatchStoreCodeProposal")(&err)

	if err := p.ValidateBasic(); err != nil {
		return err
	}

	for idx := range p.Proposals {
		if err := wasmProposalHandler(ctx, &p.Proposals[idx]); err != nil {
			return err
		}
	}

	return nil
}
