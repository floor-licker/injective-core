package ethereum

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/InjectiveLabs/metrics/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	peggytypes "github.com/InjectiveLabs/injective-core/injective-chain/modules/peggy/types"
	"github.com/InjectiveLabs/injective-core/peggo/orchestrator/ethereum/committer"
	"github.com/InjectiveLabs/injective-core/peggo/orchestrator/ethereum/peggy"
	"github.com/InjectiveLabs/injective-core/peggo/orchestrator/ethereum/provider"
	peggyevents "github.com/InjectiveLabs/injective-core/peggo/solidity/wrappers/Peggy"
)

type NetworkConfig struct {
	EthNodeRPC            string
	GasPriceAdjustment    float64
	MaxGasPrice           string
	PendingTxWaitDuration string
	EthNodeAlchemyWS      string
}

// Network is the orchestrator's reference endpoint to the Ethereum network
type Network interface {
	GetHeaderByNumber(ctx context.Context, number *big.Int) (*gethtypes.Header, error)
	GetPeggyID(ctx context.Context) (gethcommon.Hash, error)

	GetSendToInjectiveEvents(ctx context.Context,
		startBlock,
		endBlock uint64,
	) ([]*peggyevents.PeggySendToInjectiveEvent, error)
	GetPeggyERC20DeployedEvents(ctx context.Context,
		startBlock,
		endBlock uint64,
	) ([]*peggyevents.PeggyERC20DeployedEvent, error)
	GetValsetUpdatedEvents(ctx context.Context,
		startBlock,
		endBlock uint64,
	) ([]*peggyevents.PeggyValsetUpdatedEvent, error)
	GetTransactionBatchExecutedEvents(ctx context.Context,
		startBlock,
		endBlock uint64,
	) ([]*peggyevents.PeggyTransactionBatchExecutedEvent, error)

	GetValsetNonce(ctx context.Context) (*big.Int, error)
	SendEthValsetUpdate(ctx context.Context,
		oldValset *peggytypes.Valset,
		newValset *peggytypes.Valset,
		confirms []*peggytypes.MsgValsetConfirm,
	) (*gethcommon.Hash, error)

	GetTxBatchNonce(ctx context.Context, erc20ContractAddress gethcommon.Address) (*big.Int, error)
	SendTransactionBatch(ctx context.Context,
		currentValset *peggytypes.Valset,
		batch *peggytypes.OutgoingTxBatch,
		confirms []*peggytypes.MsgConfirmBatch,
	) (*gethcommon.Hash, error)

	TokenDecimals(ctx context.Context, tokenContract gethcommon.Address) (uint8, error)
}

type network struct {
	peggy.PeggyContract
	meter metrics.Meter

	FromAddr gethcommon.Address
}

func NewNetwork(
	peggyContractAddr,
	fromAddr gethcommon.Address,
	signerFn bind.SignerFn,
	cfg NetworkConfig,
	meter metrics.Meter,
) (Network, error) {
	if meter == nil {
		meter = metrics.NewNilMeter()
	}

	meter = meter.SubMeter("eth", metrics.Tag("svc", "eth"))

	evmRPC, err := rpc.Dial(cfg.EthNodeRPC)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to ethereum RPC: %s", cfg.EthNodeRPC)
	}

	ethCommitter, err := committer.NewEthCommitter(
		fromAddr,
		cfg.GasPriceAdjustment,
		cfg.MaxGasPrice,
		signerFn,
		provider.NewEVMProvider(evmRPC, meter),
		meter,
	)
	if err != nil {
		return nil, err
	}

	pendingTxDuration, err := time.ParseDuration(cfg.PendingTxWaitDuration)
	if err != nil {
		return nil, err
	}

	peggyContract, err := peggy.NewPeggyContract(ethCommitter, peggyContractAddr, peggy.PendingTxInputList{}, pendingTxDuration, meter)
	if err != nil {
		return nil, err
	}

	// If Alchemy Websocket URL is set, then Subscribe to Pending Transaction of Peggy Contract.
	if cfg.EthNodeAlchemyWS != "" {
		log.WithFields(log.Fields{
			"url": cfg.EthNodeAlchemyWS,
		}).Infoln("subscribing to Alchemy websocket")
		go peggyContract.SubscribeToPendingTxs(cfg.EthNodeAlchemyWS)
	}

	n := &network{
		PeggyContract: peggyContract,
		FromAddr:      fromAddr,
		meter:         meter,
	}

	return n, nil
}

func (n *network) TokenDecimals(ctx context.Context, tokenContract gethcommon.Address) (decimals uint8, err error) {
	ctx, done := n.meter.FuncTimingCtx(ctx, "TokenDecimals")
	defer done(&err)

	msg := ethereum.CallMsg{
		To:   &tokenContract,
		Data: gethcommon.Hex2Bytes("313ce567"), // decimals() method signature
	}

	res, err := n.Provider().CallContract(ctx, msg, nil)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, errors.Errorf("no decimals found for token contract %s", tokenContract.Hex())
	}

	return uint8(big.NewInt(0).SetBytes(res).Uint64()), nil
}

func (n *network) GetHeaderByNumber(ctx context.Context, number *big.Int) (header *gethtypes.Header, err error) {
	ctx, done := n.meter.FuncTimingCtx(ctx, "GetHeaderByNumber")
	defer done(&err)

	return n.Provider().HeaderByNumber(ctx, number)
}

func (n *network) GetPeggyID(ctx context.Context) (peggyID gethcommon.Hash, err error) {
	ctx, done := n.meter.FuncTimingCtx(ctx, "GetPeggyID")
	defer done(&err)

	return n.PeggyContract.GetPeggyID(ctx, n.FromAddr)
}

func (n *network) GetValsetNonce(ctx context.Context) (nonce *big.Int, err error) {
	ctx, done := n.meter.FuncTimingCtx(ctx, "GetValsetNonce")
	defer done(&err)

	return n.PeggyContract.GetValsetNonce(ctx, n.FromAddr)
}

func (n *network) GetTxBatchNonce(ctx context.Context, erc20ContractAddress gethcommon.Address) (nonce *big.Int, err error) {
	ctx, done := n.meter.FuncTimingCtx(ctx, "GetTxBatchNonce")
	defer done(&err)

	return n.PeggyContract.GetTxBatchNonce(ctx, erc20ContractAddress, n.FromAddr)
}

func (n *network) GetSendToInjectiveEvents(
	ctx context.Context,
	startBlock,
	endBlock uint64,
) (events []*peggyevents.PeggySendToInjectiveEvent, err error) {
	_, done := n.meter.FuncTimingCtx(ctx, "GetSendToInjectiveEvents")
	defer done(&err)

	peggyFilterer, err := peggyevents.NewPeggyFilterer(n.Address(), n.Provider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to init Peggy events filterer")
	}

	iter, err := peggyFilterer.FilterSendToInjectiveEvent(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}, nil, nil, nil)
	if err != nil {
		if !isUnknownBlockErr(err) {
			return nil, errors.Wrapf(err, "failed to scan past SendToInjectiveEvent events from Ethereum (%d - %d)", startBlock, endBlock)
		} else if iter == nil {
			return nil, errors.New("no iterator returned")
		}
	}

	defer iter.Close()

	for iter.Next() {
		events = append(events, iter.Event)
	}

	return events, nil
}

func (n *network) GetPeggyERC20DeployedEvents(
	ctx context.Context,
	startBlock,
	endBlock uint64,
) (events []*peggyevents.PeggyERC20DeployedEvent, err error) {
	_, done := n.meter.FuncTimingCtx(ctx, "GetPeggyERC20DeployedEvents")
	defer done(&err)

	peggyFilterer, err := peggyevents.NewPeggyFilterer(n.Address(), n.Provider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to init Peggy events filterer")
	}

	iter, err := peggyFilterer.FilterERC20DeployedEvent(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}, nil)
	if err != nil {
		if !isUnknownBlockErr(err) {
			return nil, errors.Wrapf(err, "failed to scan past ERC20DeployedEvent events from Ethereum (%d - %d)", startBlock, endBlock)
		} else if iter == nil {
			return nil, errors.New("no iterator returned")
		}
	}

	defer iter.Close()

	for iter.Next() {
		events = append(events, iter.Event)
	}

	return events, nil
}

func (n *network) GetValsetUpdatedEvents(ctx context.Context, startBlock, endBlock uint64) (events []*peggyevents.PeggyValsetUpdatedEvent, err error) {
	_, done := n.meter.FuncTimingCtx(ctx, "GetValsetUpdatedEvents")
	defer done(&err)

	peggyFilterer, err := peggyevents.NewPeggyFilterer(n.Address(), n.Provider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to init Peggy events filterer")
	}

	iter, err := peggyFilterer.FilterValsetUpdatedEvent(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}, nil)
	if err != nil {
		if !isUnknownBlockErr(err) {
			return nil, errors.Wrapf(err, "failed to scan past ValsetUpdatedEvent events from Ethereum (%d - %d)", startBlock, endBlock)
		} else if iter == nil {
			return nil, errors.New("no iterator returned")
		}
	}

	defer iter.Close()

	for iter.Next() {
		events = append(events, iter.Event)
	}

	return events, nil
}

func (n *network) GetTransactionBatchExecutedEvents(
	ctx context.Context,
	startBlock,
	endBlock uint64,
) (events []*peggyevents.PeggyTransactionBatchExecutedEvent, err error) {
	_, done := n.meter.FuncTimingCtx(ctx, "GetTransactionBatchExecutedEvents")
	defer done(&err)

	peggyFilterer, err := peggyevents.NewPeggyFilterer(n.Address(), n.Provider())
	if err != nil {
		return nil, errors.Wrap(err, "failed to init Peggy events filterer")
	}

	iter, err := peggyFilterer.FilterTransactionBatchExecutedEvent(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}, nil, nil)
	if err != nil {
		if !isUnknownBlockErr(err) {
			return nil, errors.Wrapf(err, "failed to scan past TransactionBatchExecuted events from Ethereum (%d - %d)", startBlock, endBlock)
		} else if iter == nil {
			return nil, errors.New("no iterator returned")
		}
	}

	defer iter.Close()

	for iter.Next() {
		events = append(events, iter.Event)
	}

	return events, nil
}

func isUnknownBlockErr(err error) bool {
	// Geth error
	if strings.Contains(err.Error(), "unknown block") {
		return true
	}

	// Parity error
	if strings.Contains(err.Error(), "One of the blocks specified in filter") {
		return true
	}

	return false
}
