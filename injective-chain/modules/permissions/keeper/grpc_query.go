package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	voucherspkg "github.com/InjectiveLabs/injective-core/injective-chain/modules/common/vouchers/types"
	"github.com/InjectiveLabs/injective-core/injective-chain/modules/permissions/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{Keeper: k}
}

func (q queryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "Params")()

	return &types.QueryParamsResponse{Params: q.GetParams(sdk.UnwrapSDKContext(c))}, nil
}

func (q queryServer) NamespaceDenoms(c context.Context, _ *types.QueryNamespaceDenomsRequest) (*types.QueryNamespaceDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "NamespaceDenoms")()

	return &types.QueryNamespaceDenomsResponse{Denoms: q.GetAllNamespaceDenoms(ctx)}, nil
}

func (q queryServer) Namespaces(c context.Context, _ *types.QueryNamespacesRequest) (*types.QueryNamespacesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "Namespaces")()

	namespaces, err := q.GetAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryNamespacesResponse{Namespaces: namespaces}, nil
}

func (q queryServer) Namespace(c context.Context, req *types.QueryNamespaceRequest) (*types.QueryNamespaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "Namespace")()

	namespace, err := q.GetNamespace(ctx, req.Denom, true)
	if err != nil {
		return nil, err
	}

	return &types.QueryNamespaceResponse{Namespace: namespace}, nil
}

func (q queryServer) RolesByActor(c context.Context, req *types.QueryRolesByActorRequest) (*types.QueryRolesByActorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "RolesByActor")()

	addr, err := sdk.AccAddressFromBech32(req.Actor)
	if err != nil {
		return nil, err
	}

	roles, err := q.GetAddressRoleNames(ctx, req.Denom, addr)
	if err != nil {
		return nil, err
	}
	return &types.QueryRolesByActorResponse{Roles: roles}, nil
}

func (q queryServer) ActorsByRole(c context.Context, req *types.QueryActorsByRoleRequest) (*types.QueryActorsByRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "ActorsByRole")()

	if !q.HasNamespace(ctx, req.Denom) {
		return nil, types.ErrUnknownDenom
	}

	role, err := q.GetRoleByName(ctx, req.Denom, req.Role)
	if err != nil {
		return nil, err
	}

	actors := make([]string, 0)

	// inefficient but the only way for now since we don't index actors by roleID
	err = q.IterateActorRoles(ctx, req.Denom, func(actor sdk.AccAddress, roleIDs []uint32) error {
		for _, roleID := range roleIDs {
			if roleID != role.RoleId {
				continue
			}
			actors = append(actors, actor.String())
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryActorsByRoleResponse{
		Actors: actors,
	}, nil
}

func (q queryServer) RoleManagers(c context.Context, req *types.QueryRoleManagersRequest) (*types.QueryRoleManagersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "RoleManagers")()

	if !q.HasNamespace(ctx, req.Denom) {
		return nil, types.ErrUnknownDenom
	}

	managers, err := q.GetAllRoleManagers(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryRoleManagersResponse{RoleManagers: managers}, nil
}

func (q queryServer) RoleManager(c context.Context, req *types.QueryRoleManagerRequest) (*types.QueryRoleManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "RoleManager")()

	if !q.HasNamespace(ctx, req.Denom) {
		return nil, types.ErrUnknownDenom
	}

	manager, err := sdk.AccAddressFromBech32(req.Manager)
	if err != nil {
		return nil, err
	}

	roleIDs := q.getAllRolesIDsForManager(ctx, req.Denom, manager)

	roles := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		role, err := q.GetRoleByID(ctx, req.Denom, roleID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role.Name)
	}
	roleManager := types.NewRoleManager(manager, roles)

	return &types.QueryRoleManagerResponse{RoleManager: roleManager}, nil
}

func (q queryServer) PolicyStatuses(c context.Context, req *types.QueryPolicyStatusesRequest) (*types.QueryPolicyStatusesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "PolicyStatuses")()

	if !q.HasNamespace(ctx, req.Denom) {
		return nil, types.ErrUnknownDenom
	}

	statuses, err := q.GetAllPolicyStatuses(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryPolicyStatusesResponse{
		PolicyStatuses: statuses,
	}, nil
}

func (q queryServer) PolicyManagerCapabilities(c context.Context, req *types.QueryPolicyManagerCapabilitiesRequest) (*types.QueryPolicyManagerCapabilitiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "PolicyManagerCapabilities")()

	if !q.HasNamespace(ctx, req.Denom) {
		return nil, types.ErrUnknownDenom
	}

	capabilities, err := q.GetAllPolicyManagerCapabilities(ctx, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryPolicyManagerCapabilitiesResponse{
		PolicyManagerCapabilities: capabilities,
	}, nil
}

func (q queryServer) Vouchers(c context.Context, req *types.QueryVouchersRequest) (res *types.QueryVouchersResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "Vouchers")(&err)

	var raw []voucherspkg.AddressVoucher
	if req.Denom == "" {
		raw, err = q.vouchersAssistant.GetAllVouchers(ctx)
	} else {
		raw, err = q.vouchersAssistant.GetVouchersForDenom(ctx, req.Denom)
	}
	if err != nil {
		return nil, err
	}

	res = &types.QueryVouchersResponse{Vouchers: raw}
	return res, nil
}

func (q queryServer) Voucher(c context.Context, req *types.QueryVoucherRequest) (res *types.QueryVoucherResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "Voucher")(&err)

	if req.Denom == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "denom is required")
	}
	if req.Address == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "address is required")
	}

	var addr sdk.AccAddress

	addr, err = sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	var voucher sdk.Coin

	voucher, err = q.vouchersAssistant.GetVoucher(ctx, req.Denom, addr)
	if err != nil {
		return nil, err
	}

	res = &types.QueryVoucherResponse{Voucher: voucher}
	return res, nil
}

func (q queryServer) PermissionsModuleState(c context.Context, req *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer q.Meter(ctx).FuncTiming(&ctx, "PermissionsModuleState")()

	return &types.QueryModuleStateResponse{State: q.ExportGenesis(ctx)}, nil
}
