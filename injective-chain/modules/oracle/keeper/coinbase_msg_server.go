package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/oracle/types"
)

type CoinbaseMsgServer struct {
	*Keeper
}

// NewCoinbaseMsgServerImpl returns an implementation of the coinbase provider MsgServer interface for the provided Keeper for coinbase provider oracle functions.
func NewCoinbaseMsgServerImpl(keeper Keeper) CoinbaseMsgServer {
	return CoinbaseMsgServer{
		Keeper: &keeper,
	}
}

func (k CoinbaseMsgServer) RelayCoinbaseMessages(c context.Context, msg *types.MsgRelayCoinbaseMessages) (*types.MsgRelayCoinbaseMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	defer k.Meter(ctx).FuncTiming(&ctx, "RelayCoinbaseMessages")()

	for idx := range msg.Messages {
		err := types.ValidateCoinbaseSignature(msg.Messages[idx], msg.Signatures[idx])
		if err != nil {
			return nil, err
		}

		newCoinbasePriceState, err := types.ParseCoinbaseMessage(msg.Messages[idx])
		if err != nil {
			return nil, err
		}

		price := newCoinbasePriceState.GetDecPrice()

		oldCoinbasePriceState := k.getLastCoinbasePriceState(ctx, newCoinbasePriceState.Key)
		blockTime := ctx.BlockTime().Unix()
		if oldCoinbasePriceState == nil {
			newCoinbasePriceState.PriceState = types.PriceState{
				Price:           price,
				CumulativePrice: math.LegacyZeroDec(),
				Timestamp:       blockTime,
			}
		} else {
			oldCoinbasePriceState.PriceState.UpdatePrice(price, blockTime)
			newCoinbasePriceState.PriceState = oldCoinbasePriceState.PriceState
		}

		if err = k.SetCoinbasePriceState(ctx, newCoinbasePriceState); err != nil {
			return nil, err
		}
	}

	return &types.MsgRelayCoinbaseMessagesResponse{}, nil
}
