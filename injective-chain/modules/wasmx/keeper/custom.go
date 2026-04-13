package keeper

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/types"
)

func (k *Keeper) InjectiveExec(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	funds sdk.Coins,
	msg *types.InjectiveExecMsg,
) ([]byte, error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "InjectiveExec")()

	execBz, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := k.wasmContractOpsKeeper.Execute(
		ctx,
		contractAddress,
		contractAddress,
		execBz,
		funds,
	)
	if err != nil {
		k.Logger(ctx).Debug("result", res, "err", err)
		return res, err
	}

	k.Logger(ctx).Debug("InjectiveExec result:", "result", string(res))
	return res, nil
}

func (k *Keeper) PinContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "PinContract")()

	contractInfo := k.wasmViewKeeper.GetContractInfo(ctx, contractAddress)
	if contractInfo == nil {
		return errors.Wrapf(
			sdkerrors.ErrNotFound,
			"Contract with address %v not found",
			contractAddress.String(),
		)
	}
	err = k.wasmContractOpsKeeper.PinCode(ctx, contractInfo.CodeID)
	if err != nil {
		// Wasmer runtime error
		k.Logger(ctx).
			Error("❌ Error while pinning the contract", "contractAddress", contractAddress.String(), "error", err)
		return
	}

	k.Logger(ctx).
		Debug("✅ Pinned the contract successfully", "contractAddress", contractAddress.String())
	return nil
}

func (k *Keeper) UnpinContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
) (err error) {
	defer k.Meter(ctx).FuncTiming(&ctx, "UnpinContract")()

	contractInfo := k.wasmViewKeeper.GetContractInfo(ctx, contractAddress)
	if contractInfo == nil {
		return errors.Wrapf(
			sdkerrors.ErrNotFound,
			"Contract with address %v not found",
			contractAddress.String(),
		)
	}
	err = k.wasmContractOpsKeeper.UnpinCode(ctx, contractInfo.CodeID)
	if err != nil {
		// Wasmer runtime error
		k.Logger(ctx).
			Error("❌ Error while unpinning the contract", "contractAddress", contractAddress.String(), "error", err)
		return
	}

	k.Logger(ctx).
		Debug("✅ Unpinned the contract successfully", "contractAddress", contractAddress.String())
	return nil
}
