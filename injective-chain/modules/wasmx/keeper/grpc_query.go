package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/wasmx/types"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) WasmxParams(c context.Context, _ *types.QueryWasmxParamsRequest) (*types.QueryWasmxParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "WasmxParams")()

	params := k.GetParams(ctx)

	res := &types.QueryWasmxParamsResponse{
		Params: params,
	}
	return res, nil
}

func (k *Keeper) ContractRegistrationInfo(c context.Context, req *types.QueryContractRegistrationInfoRequest) (*types.QueryContractRegistrationInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "ContractRegistrationInfo")()

	contract, err := sdk.AccAddressFromBech32(req.ContractAddress)

	if err != nil {
		return nil, types.ErrInvalidContractAddress
	}

	res := &types.QueryContractRegistrationInfoResponse{
		Contract: k.GetContractByAddress(ctx, contract),
	}
	return res, nil
}

func (k *Keeper) WasmxModuleState(c context.Context, _ *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "WasmxModuleState")()

	res := &types.QueryModuleStateResponse{
		State: k.ExportGenesis(ctx),
	}
	return res, nil
}
