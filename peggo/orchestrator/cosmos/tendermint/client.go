package tendermint

import (
	"context"

	"github.com/InjectiveLabs/metrics/v2"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	comettypes "github.com/cometbft/cometbft/rpc/core/types"
	log "github.com/xlab/suplog"
)

type Client interface {
	GetBlock(ctx context.Context, height int64) (*comettypes.ResultBlock, error)
}

type tmClient struct {
	rpcClient rpcclient.Client
	meter     metrics.Meter
}

func NewRPCClient(rpcNodeAddr string, meter metrics.Meter) Client {
	if meter == nil {
		meter = metrics.NewNilMeter()
	}

	rpcClient, err := rpchttp.NewWithTimeout(rpcNodeAddr, 10)
	if err != nil {
		log.WithError(err).Fatalln("failed to init rpcClient")
	}

	return &tmClient{
		rpcClient: rpcClient,
		meter:     meter.SubMeter("tendermint", metrics.Tag("svc", "tendermint")),
	}
}

// GetBlock queries for a block by height
func (c *tmClient) GetBlock(ctx context.Context, height int64) (block *comettypes.ResultBlock, err error) {
	ctx, done := c.meter.FuncTimingCtx(ctx, "GetBlock")
	defer done(&err)

	return c.rpcClient.Block(ctx, &height)
}
