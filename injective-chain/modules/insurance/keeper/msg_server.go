package keeper

import (
	"context"

	"cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/insurance/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the insurance MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) UpdateParams(c context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "UpdateParams")()

	if msg.Authority != k.authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority: expected %s, got %s", k.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	k.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

// CreateInsuranceFund is wrapper of keeper.CreateInsuranceFund
func (k msgServer) CreateInsuranceFund(c context.Context, msg *types.MsgCreateInsuranceFund) (*types.MsgCreateInsuranceFundResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "CreateInsuranceFund")()

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	isPerpetualOrBinaryOptionsExpirationFlag := msg.Expiry == types.PerpetualExpiryFlag || msg.Expiry == types.BinaryOptionsExpiryFlag
	if !isPerpetualOrBinaryOptionsExpirationFlag && msg.Expiry < ctx.BlockTime().Unix() {
		return nil, types.ErrInvalidExpirationTime
	}

	if err := k.Keeper.CreateInsuranceFund(ctx, sender, msg.InitialDeposit, msg.Ticker, msg.QuoteDenom, msg.OracleBase, msg.OracleQuote, msg.OracleType, msg.Expiry); err != nil {
		k.Logger(ctx).Error("Insurance fund creation failed", err)
		return nil, err
	}

	return &types.MsgCreateInsuranceFundResponse{}, nil
}

// Underwrite is wrapper of keeper.UnderwriteInsuranceFund
func (k msgServer) Underwrite(c context.Context, msg *types.MsgUnderwrite) (*types.MsgUnderwriteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "Underwrite")()

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	marketID := common.HexToHash(msg.MarketId)
	if err := k.Keeper.UnderwriteInsuranceFund(ctx, sender, marketID, msg.Deposit); err != nil {
		k.Logger(ctx).Error("underwriting insurance fund failed", err)
		return nil, err
	}

	return &types.MsgUnderwriteResponse{}, nil
}

// RequestRedemption is wrapper of keeper.RequestInsuranceFundRedemption
func (k msgServer) RequestRedemption(c context.Context, msg *types.MsgRequestRedemption) (*types.MsgRequestRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RequestRedemption")()

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	marketID := common.HexToHash(msg.MarketId)
	err = k.Keeper.RequestInsuranceFundRedemption(ctx, sender, marketID, msg.Amount)
	if err != nil {
		k.Logger(ctx).Error("requesting redemption for insurance fund failed", err)
		return nil, err
	}

	return &types.MsgRequestRedemptionResponse{}, nil
}

func (k msgServer) ClaimVoucher(goCtx context.Context, msg *types.MsgClaimVoucher) (res *types.MsgClaimVoucherResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	defer k.Meter(ctx).FuncTiming(&ctx, "ClaimVoucher")(&err)

	senderAddr := sdk.MustAccAddressFromBech32(msg.Sender)

	err = k.vouchersAssistant.ClaimVoucher(ctx, senderAddr, msg.Denom)
	if err != nil {
		return nil, err
	}

	res = &types.MsgClaimVoucherResponse{}
	return res, nil
}
