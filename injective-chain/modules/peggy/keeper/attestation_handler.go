package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/modules/peggy/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper     *Keeper
	bankKeeper types.BankKeeper
}

func NewAttestationHandler(bankKeeper types.BankKeeper, keeper Keeper) AttestationHandler {
	return AttestationHandler{
		keeper:     &keeper,
		bankKeeper: bankKeeper,
	}
}

// Handle is the entry point for Attestation processing.
func (h AttestationHandler) Handle(ctx sdk.Context, claim types.EthereumClaim) error {
	defer h.keeper.Meter(ctx).FuncTiming(&ctx, "AttestationHandler.Handle")()

	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		return h.handleDepositClaim(ctx, claim)
	case *types.MsgWithdrawClaim:
		h.handleWithdrawClaim(ctx, claim)
	case *types.MsgERC20DeployedClaim:
		// todo: upgrade Peggy.sol on testnet
		// The deployERC20 functionality was removed from mainnet contract.
		// Logic below no longer applies so we return early with no error
	case *types.MsgValsetUpdatedClaim:
		h.handleValsetUpdatedClaim(ctx, claim)
	default:
		return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Invalid event type for attestations %s", claim.GetType()))
	}

	return nil
}

func (h AttestationHandler) handleDepositClaim(ctx sdk.Context, claim *types.MsgDepositClaim) error {
	defer h.keeper.Meter(ctx).FuncTiming(&ctx, "AttestationHandler.handleDepositClaim")()

	sender, err := types.NewEthAddress(claim.EthereumSender)
	if err != nil {
		// likewise nil sender would have to be caused by a bogus event
		return errors.Wrap(err, "failed to parse ethereum sender in claim")
	}

	// Check if coin is Cosmos-originated asset and get denom
	tokenContract := common.HexToAddress(claim.TokenContract)
	isCosmosOriginated, denom := h.keeper.ERC20ToDenomLookup(ctx, tokenContract)
	depositCoin := sdk.NewCoin(denom, claim.Amount)

	if !isCosmosOriginated {
		// Check if supply overflows with claim amount (first pass)
		currentSupply := h.bankKeeper.GetSupply(ctx, denom)
		newSupply := new(big.Int).Add(currentSupply.Amount.BigInt(), claim.Amount.BigInt())
		if newSupply.BitLen() > 256 {
			return errors.Wrap(types.ErrSupplyOverflow, "invalid coin supply")
		}

		if err := h.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(depositCoin)); err != nil {
			return errors.Wrapf(err, "failed to mint deposit coin: %s", depositCoin.String())
		}

		// Interestingly, even tho INJ is cosmos native by nature, it is interpreted as an ERC20.
		// We only track mint amounts for assets that are truly eth native.
		if denom != chaintypes.InjectiveCoin {
			// increment tracked mint amount (2nd pass)
			// because we burn bank on withdrawal inclusion not execution
			currentAmount := h.keeper.GetMintAmountERC20(ctx, tokenContract)
			newAmount, err := currentAmount.SafeAdd(claim.Amount)
			if err != nil {
				// what this here means is that someone tried to bridge in a custom ERC20 token with malicious behavior
				// one that allows depositing amounts larger than possibly representable on a cosmos chain.
				// Event execution becomes a no-op except for the nonce increment, and we move on (like with all attestation errors)
				return errors.Wrapf(ErrAbsoluteMintLimitOverflow, "failed to mint coin amount: %v", err)
			}

			h.keeper.SetMintAmountERC20(ctx, tokenContract, newAmount)
		}
	}

	receiver := sdk.MustAccAddressFromBech32(claim.CosmosReceiver)
	if h.keeper.IsOnBlacklist(ctx, *sender) {
		// sender is blacklisted, we deposit to segregated wallet
		receiver = sdk.MustAccAddressFromBech32(h.keeper.GetParams(ctx).SegregatedWalletAddress)
		if err := h.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(depositCoin)); err != nil {
			return errors.Wrap(err, "failed to send sanctioned deposit to segregated wallet")
		}

		h.keeper.TrackTokenInflow(ctx, tokenContract, depositCoin.Amount)
		_ = ctx.EventManager().EmitTypedEvent(types.NewEventDepositReceived(*sender, receiver, depositCoin))
		return nil
	}

	// address appears valid, attempt to send minted/locked coins to receiver
	if err := h.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(depositCoin)); err != nil {
		// last attempt, send to community pool (default behavior)
		if err := h.keeper.SendToCommunityPool(ctx, sdk.NewCoins(depositCoin)); err != nil {
			return errors.Wrap(err, "failed to send deposit to community pool")
		}

		receiver = h.keeper.accountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName).GetAddress()
	}

	h.keeper.TrackTokenInflow(ctx, tokenContract, depositCoin.Amount)

	_ = ctx.EventManager().EmitTypedEvent(types.NewEventDepositReceived(*sender, receiver, depositCoin))
	return nil
}

func (h AttestationHandler) handleWithdrawClaim(ctx sdk.Context, claim *types.MsgWithdrawClaim) {
	defer h.keeper.Meter(ctx).FuncTiming(&ctx, "AttestationHandler.handleWithdrawClaim")()

	h.keeper.OutgoingTxBatchExecuted(ctx, common.HexToAddress(claim.TokenContract), claim.BatchNonce)
}

func (h AttestationHandler) handleValsetUpdatedClaim(ctx sdk.Context, claim *types.MsgValsetUpdatedClaim) {
	defer h.keeper.Meter(ctx).FuncTiming(&ctx, "AttestationHandler.handleValsetUpdatedClaim")()
	// TODO here we should check the contents of the validator set against
	// the store, if they differ we should take some action to indicate to the
	// user that bridge highjacking has occurred
	h.keeper.SetLastObservedValset(ctx, types.Valset{
		Nonce:        claim.ValsetNonce,
		Members:      claim.Members,
		RewardAmount: claim.RewardAmount,
		RewardToken:  claim.RewardToken,
	})
}
