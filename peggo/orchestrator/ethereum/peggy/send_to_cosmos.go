package peggy

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	wrappers "github.com/InjectiveLabs/injective-core/peggo/solidity/wrappers/InjToken"
)

func (s *peggyContract) SendToCosmos(
	ctx context.Context,
	erc20 common.Address,
	amount *big.Int,
	cosmosAccAddress sdk.AccAddress,
	senderAddress common.Address,
) (txHash *common.Hash, err error) {
	ctx, done := s.meter.FuncTimingCtx(ctx, "SendToCosmos")
	defer done(&err)

	erc20Wrapper, err := wrappers.NewERC20(erc20, s.ethProvider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ERC20 wrapper")
	}

	if allowance, err := erc20Wrapper.Allowance(&bind.CallOpts{
		From:    common.Address{},
		Context: ctx,
	}, senderAddress, s.peggyAddress); err != nil {
		return nil, errors.Wrap(err, "failed to get ERC20 allowance for peggy contract")
	} else if allowance.Cmp(maxUintAllowance) != 0 {
		// allowance not set or not max (a.k.a. unlocked token)
		txData, err := erc20ABI.Pack("approve", s.peggyAddress, maxUintAllowance)
		if err != nil {
			log.WithError(err).Errorln("ABI Pack (ERC20 approve) method")
			return nil, err
		}

		hash, err := s.SendTx(ctx, erc20, txData)
		if err != nil {
			log.WithError(err).WithField("tx_hash", hash.Hex()).Errorln("Failed to sign and submit (ERC20 approve) to EVM")
			return nil, err
		}

		log.Infoln("Sent Tx (ERC20 approve):", hash.Hex())
	}

	// This code deals with some specifics of Ethereum byte encoding, Ethereum is BigEndian
	// so small values like addresses that don't take up the full length of the byte vector
	// are pushed up to the top. This duplicates the way Ethereum encodes it's own addresses
	// as closely as possible.
	cosmosDestAddressBytes := cosmosAccAddress.Bytes()
	for len(cosmosDestAddressBytes) < 32 {
		cosmosDestAddressBytes = append(cosmosDestAddressBytes, byte(0))
	}

	txData, err := peggyABI.Pack("sendToCosmos", erc20, cosmosDestAddressBytes, amount)
	if err != nil {
		log.WithError(err).Errorln("ABI Pack (Peggy sendToCosmos) method")
		return nil, err
	}

	hash, err := s.SendTx(ctx, s.peggyAddress, txData)
	if err != nil {
		log.WithError(err).WithField("tx_hash", hash.Hex()).Errorln("Failed to sign and submit (Peggy sendToCosmos) to EVM")
		return nil, err
	}

	log.Infoln("Sent Tx (Peggy sendToCosmos):", hash.Hex())

	return &hash, nil
}
