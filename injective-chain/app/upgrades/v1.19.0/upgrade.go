package v1dot19dot0

import (
	"bytes"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/ethereum/go-ethereum/common"

	"github.com/InjectiveLabs/injective-core/injective-chain/app/upgrades"
	exchangetypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types"
	exchangev2types "github.com/InjectiveLabs/injective-core/injective-chain/modules/exchange/types/v2"
	peggytypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/peggy/types"
	chaintypes "github.com/InjectiveLabs/injective-core/injective-chain/types"
)

const (
	UpgradeVersion = "v1.19.0"
)

var ausdSpotMarketIDsToDemolish = []common.Hash{
	common.HexToHash("0xae91b4f0c34ed3e7fbd9fbec1a9acbc74d53892de5e79e8da8f9537471dbd6b7"), // ausd-usdt
	common.HexToHash("0x609376b88fd2e6e37ffc2ac9d6858633ab1d6c90ec6d100f9d205d38e531f8ac"), // ausd-usdc
	common.HexToHash("0xee25f87facc4952d06b93d44ebb0e9eebf9ee1b39c6c7dcf6b14445b6bb6097e"), // inj-ausd
}

var whiteKnightLiquidatorsToAdd = []string{
	"inj1jr094svc02v8ft5x3ecrcnrryp2g5vt23ngulu",
	"inj1fsuxt5a0cfsxpj6gzxx5nmyhd9hy97edklj4cn",
	"inj16ezzzulhz66v5yjh0nkuutms0at66dvh8d40ew",
	"inj1xdqvnlk3jztnfc79mula20cvlhagtauneslm96",
	"inj1qqyvtnzpygeqe8d9y3ghmsj020ljdly6nc9ght",
	"inj1h2ltgg6kk237f53f3fgph0tc6w5lqwv2whemmz",
	"inj12fgfe6mkvv7rfslcmj86nnvvajsknu45lkpkjk",
	"inj174splsttpxxrrwcew0nwxy83nf2tmx2nvkjjru",
}
func StoreUpgrades() storetypes.StoreUpgrades {
	return storetypes.StoreUpgrades{
		Added:   nil,
		Renamed: nil,
		Deleted: []string{"chainlink"},
	}
}

func UpgradeSteps() []*upgrades.UpgradeHandlerStep {
	return []*upgrades.UpgradeHandlerStep{
		upgrades.NewUpgradeHandlerStep(
			"[peggy] Migrate MintAmountERC20",
			UpgradeVersion,
			upgrades.MainnetChainID,
			MigrateMintAmountERC20,
		),
		upgrades.NewUpgradeHandlerStep(
			"[peggy] Update MsgValsetUpdatedClaim entries with new key",
			UpgradeVersion,
			upgrades.MainnetChainID,
			UpdatePeggyClaimHashEntries,
		),
		upgrades.NewUpgradeHandlerStep(
			"[exchange] Backfill insurance redemption addr index",
			UpgradeVersion,
			upgrades.MainnetChainID,
			BackfillInsuranceRedemptionAddrIndex,
		),
		upgrades.NewUpgradeHandlerStep(
			"[exchange] Create auction fees subaccount",
			UpgradeVersion,
			upgrades.MainnetChainID,
			CreateAuctionFeesSubaccount,
		),
		upgrades.NewUpgradeHandlerStep(
			"[exchange] Demolish legacy aUSD spot markets",
			UpgradeVersion,
			upgrades.MainnetChainID,
			DemolishAUSDSpotMarkets,
		),
		upgrades.NewUpgradeHandlerStep(
			"[exchange] Update auction InjBasketMaxCap",
			UpgradeVersion,
			upgrades.MainnetChainID,
			UpdateAuctionInjBasketMaxCap,
		),
		upgrades.NewUpgradeHandlerStep(
			"[authz] Rebuild x/authz grants to create new indexes for USDC Blocklisted event listener",
			UpgradeVersion,
			upgrades.MainnetChainID,
			RebuildAuthzGrantsToCreateIndexes,
		),
		upgrades.NewUpgradeHandlerStep(
			"[exchange] Configure white-knight liquidator params",
			UpgradeVersion,
			upgrades.MainnetChainID,
			ConfigureWhiteKnightLiquidatorParams,
		),
	}
}

func MigrateMintAmountERC20(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	// 0. Zero out all keys in MintAmounts
	peggyStore := ctx.KVStore(app.GetKey(peggytypes.StoreKey))
	mintAmountStore := prefix.NewStore(peggyStore, peggytypes.MintAmountERC20Key)
	keysToRemove := make([][]byte, 0)
	chaintypes.IterateKeysSafe(mintAmountStore.Iterator(nil, nil), func(k []byte) (stop bool) {
		keysToRemove = append(keysToRemove, bytes.Clone(k))
		return false
	})

	for _, key := range keysToRemove {
		mintAmountStore.Delete(key)
	}

	// 1. Collect existing amounts that have not yet been withdrawn
	peggyCoins := make(sdk.Coins, 0)
	collect := func(coin sdk.Coin) bool {
		// inj can't be picked up here anyway
		if _, err := peggytypes.NewPeggyDenomFromString(coin.Denom); err == nil {
			peggyCoins = append(peggyCoins, coin)
		}

		return false
	}

	app.GetBankKeeper().IterateTotalSupply(ctx, collect)
	pk := app.GetPeggyKeeper()

	for _, coin := range peggyCoins {
		denom, _ := peggytypes.NewPeggyDenomFromString(coin.Denom)
		token, _ := denom.TokenContract()
		pk.SetMintAmountERC20(ctx, token, coin.Amount)
	}

	// 2. If there were some withdrawals and batches present in the state,
	// they were already burned in bank during the initial send. We account for those too
	for _, withdrawal := range pk.GetPoolTransactions(ctx) {
		token := common.HexToAddress(withdrawal.Erc20Token.Contract)
		if cosmosOriginated, denom := pk.ERC20ToDenomLookup(ctx, token); !cosmosOriginated && denom != chaintypes.InjectiveCoin {
			amountToAdd := withdrawal.Erc20Token.Amount.Add(withdrawal.Erc20Fee.Amount)
			pk.SetMintAmountERC20(ctx, token, pk.GetMintAmountERC20(ctx, token).Add(amountToAdd))
		}
	}

	for _, batch := range pk.GetOutgoingTxBatches(ctx) {
		token := common.HexToAddress(batch.TokenContract)
		amountToAdd := math.ZeroInt()
		if cosmosOriginated, denom := pk.ERC20ToDenomLookup(ctx, token); !cosmosOriginated && denom != chaintypes.InjectiveCoin {
			for _, tx := range batch.Transactions {
				amountToAdd = amountToAdd.Add(tx.Erc20Token.Amount)
				amountToAdd = amountToAdd.Add(tx.Erc20Fee.Amount)
			}

			pk.SetMintAmountERC20(ctx, token, pk.GetMintAmountERC20(ctx, token).Add(amountToAdd))
		}
	}

	return nil
}

// BackfillInsuranceRedemptionAddrIndex populates the secondary (redeemer, marketID)
// index for all existing redemption schedules so that the new batching logic in
// RequestInsuranceFundRedemption can efficiently look up schedules by address.
func BackfillInsuranceRedemptionAddrIndex(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	k := app.GetInsuranceKeeper()
	k.BackfillRedemptionScheduleAddrIndex(ctx)
	logger.Info("Backfilled insurance redemption schedule address index")
	return nil
}

func DemolishAUSDSpotMarkets(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	exchangeKeeper := app.GetExchangeKeeper()

	for _, marketID := range ausdSpotMarketIDsToDemolish {
		market := exchangeKeeper.GetSpotMarketByID(ctx, marketID)
		if market == nil {
			logger.Error("Skipping missing aUSD spot market demolition target", "market_id", marketID.Hex())
			continue
		}

		if market.Status == exchangev2types.MarketStatus_Demolished {
			logger.Info("aUSD spot market already demolished", "market_id", marketID.Hex(), "ticker", market.Ticker)
			continue
		}

		exchangeKeeper.CancelAllRestingLimitOrdersFromSpotMarket(ctx, market, marketID)
		if _, err := exchangeKeeper.SetSpotMarketStatus(ctx, marketID, exchangev2types.MarketStatus_Demolished); err != nil {
			return err
		}

		logger.Info("Demolished aUSD spot market", "market_id", marketID.Hex(), "ticker", market.Ticker)
	}

	return nil
}

func UpdateAuctionInjBasketMaxCap(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	auctionKeeper := app.GetAuctionKeeper()
	params := auctionKeeper.GetParams(ctx)
	params.InjBasketMaxCap = math.NewIntWithDecimal(1_000_000, 18)
	auctionKeeper.SetParams(ctx, params)
	logger.Info("Updated auction InjBasketMaxCap", "new_value", params.InjBasketMaxCap.String())
	return nil
}

func ConfigureWhiteKnightLiquidatorParams(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	params := app.GetExchangeKeeper().GetParams(ctx)
	params.WhiteKnightLiquidatorRewardShareRate = math.LegacyNewDecWithPrec(5, 1)
	if params.LiquidatorRewardShareRate.GT(params.WhiteKnightLiquidatorRewardShareRate) {
		params.WhiteKnightLiquidatorRewardShareRate = params.LiquidatorRewardShareRate
	}

	if err := exchangetypes.ValidateWhiteKnightLiquidators(whiteKnightLiquidatorsToAdd); err != nil {
		return fmt.Errorf("invalid white-knight liquidators: %w", err)
	}

	params.WhiteKnightLiquidators = append([]string(nil), whiteKnightLiquidatorsToAdd...)
	app.GetExchangeKeeper().SetParams(ctx, params)

	logger.Info(
		"Configured white-knight liquidator params",
		"reward_share_rate",
		params.WhiteKnightLiquidatorRewardShareRate.String(),
		"num_liquidators",
		len(params.WhiteKnightLiquidators),
	)

	return nil
}

// CreateAuctionFeesSubaccount creates the auction module's fee collector subaccount (AuctionFeesSubaccountAddress)
// if it does not exist. Required for existing chains upgrading to the version where ante fees are sent to this
// subaccount instead of the main auction module; new chains get it from InitGenesis via CreateModuleAccount.
func CreateAuctionFeesSubaccount(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	app.GetAuctionKeeper().EnsureAuctionFeesSubaccount(ctx)
	logger.Info("Ensured auction fees subaccount exists")

	return nil
}

func UpdatePeggyClaimHashEntries(ctx sdk.Context, app upgrades.InjectiveApplication, _ log.Logger) error {
	// A peggy claim hash is used to uniquely identify an attestation in the state where votes aggregate.
	// Since we changed the hashing method of MsgValsetUpdatedClaim in this upgrade, old entries need to be
	// taken care of because we lost the ability to construct the old key under which the entry is stored.

	// This upgrade handler is meant to prevent the following potential scenario:
	// 	1. There's an existing MsgValsetUpdatedClaim before this upgrade
	// 	2. Some of the validators already attested to this claim (under the old claim hash key)
	//	3. In case we don't have this handler, then subsequent claims for the same event will create a separate attestation (under the new key)
	// 	4. If the timing of the upgrade was bad, we might have 2 different attestations for the same claim with no majority in any

	var (
		pk           = app.GetPeggyKeeper()
		oldKeys      = make([][]byte, 0)
		claims       = make([]peggytypes.EthereumClaim, 0)
		attestations = make([]*peggytypes.Attestation, 0)
	)

	// 1. check if there are MsgValsetUpdatedClaim attestations
	var unpackErr error
	pk.IterateAttestations(ctx, func(k []byte, v *peggytypes.Attestation) (stop bool) {
		claim, err := pk.UnpackAttestationClaim(v)
		if err != nil {
			unpackErr = err
			return true
		}

		if _, ok := claim.(*peggytypes.MsgValsetUpdatedClaim); !ok {
			return false // skip
		}

		claims = append(claims, claim)
		oldKeys = append(oldKeys, bytes.Clone(k))
		attestations = append(attestations, v)

		return false
	})

	if unpackErr != nil {
		return fmt.Errorf("failed to unpack existing attestation claim: %w", unpackErr)
	}

	if len(oldKeys) == 0 {
		return nil // no-op
	}

	// 2. remove entries under the old keys
	peggyStore := ctx.KVStore(app.GetKey(peggytypes.StoreKey))
	for _, key := range oldKeys {
		peggyStore.Delete(key)
	}

	// 3. index previous entries with the new key
	for i := 0; i < len(attestations); i++ {
		pk.SetAttestation(ctx, claims[i].GetEventNonce(), claims[i].ClaimHash(), attestations[i])
	}

	return nil
}

func RebuildAuthzGrantsToCreateIndexes(ctx sdk.Context, app upgrades.InjectiveApplication, logger log.Logger) error {
	var (
		granters []sdk.AccAddress
		grantees []sdk.AccAddress
		grants   []authz.Grant
	)

	authzKeeper := app.GetAuthzKeeper()

	// sanity check - delete grants that went expired during upgrade downtime
	err := authzKeeper.DequeueAndDeleteExpiredGrants(ctx)
	if err != nil {
		return err
	}

	authzKeeper.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant authz.Grant) bool {
		granters = append(granters, granter)
		grantees = append(grantees, grantee)
		grants = append(grants, grant)
		return false
	})

	for i := range granters {
		auth, err := grants[i].GetAuthorization()
		if err != nil {
			return fmt.Errorf("can't get Grant Authorization: %w", err)
		}
		err = authzKeeper.DeleteGrant(ctx, grantees[i], granters[i], auth.MsgTypeURL())
		if err != nil {
			return fmt.Errorf("can't delete Grant: %w", err)
		}
		err = authzKeeper.SaveGrant(ctx, grantees[i], granters[i], auth, grants[i].Expiration)
		if err != nil {
			return fmt.Errorf("can't save Grant: %w", err)
		}
	}

	logger.Info("Re-created authz Grants", "total", len(granters))

	return nil
}
