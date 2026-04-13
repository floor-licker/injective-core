package main

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	FlagJSONRPCEnable              = "json-rpc.enable"
	FlagJSONRPCAPI                 = "json-rpc.api"
	FlagJSONRPCAddress             = "json-rpc.address"
	FlagJSONWsAddress              = "json-rpc.ws-address"
	FlagJSONRPCGasCap              = "json-rpc.gas-cap"
	FlagJSONRPCEVMTimeout          = "json-rpc.evm-timeout"
	FlagJSONRPCTxFeeCap            = "json-rpc.txfee-cap"
	FlagJSONRPCFilterCap           = "json-rpc.filter-cap"
	FlagJSONRPCFeeHistoryCap       = "json-rpc.feehistory-cap"
	FlagJSONRPCLogsCap             = "json-rpc.logs-cap"
	FlagJSONRPCBlockRangeCap       = "json-rpc.block-range-cap"
	FlagJSONRPCHTTPTimeout         = "json-rpc.http-timeout"
	FlagJSONRPCHTTPIdleTimeout     = "json-rpc.http-idle-timeout"
	FlagJSONRPCAllowUnprotectedTxs = "json-rpc.allow-unprotected-txs"
	FlagJSONRPCMaxOpenConnections  = "json-rpc.max-open-connections"
	FlagJSONRPCEnableIndexer       = "json-rpc.enable-indexer"
	FlagJSONRPCAllowIndexerGap     = "json-rpc.allow-indexer-gap"
	FlagJSONRPCEnableMetrics       = "json-rpc.metrics"
	FlagJSONRPCMetricsAddress      = "json-rpc.metrics-address"
	FlagJSONRPCReturnDataLimit     = "json-rpc.return-data-limit"

	FlagJSONRPCDebugEnable             = "json-rpc-debug.enable"
	FlagJSONRPCDebugAPI                = "json-rpc-debug.api"
	FlagJSONRPCDebugAddress            = "json-rpc-debug.address"
	FlagJSONRPCDebugGasCap             = "json-rpc-debug.gas-cap"
	FlagJSONRPCDebugEVMTimeout         = "json-rpc-debug.evm-timeout"
	FlagJSONRPCDebugTxFeeCap           = "json-rpc-debug.txfee-cap"
	FlagJSONRPCDebugFilterCap          = "json-rpc-debug.filter-cap"
	FlagJSONRPCDebugFeeHistoryCap      = "json-rpc-debug.feehistory-cap"
	FlagJSONRPCDebugLogsCap            = "json-rpc-debug.logs-cap"
	FlagJSONRPCDebugBlockRangeCap      = "json-rpc-debug.block-range-cap"
	FlagJSONRPCDebugHTTPTimeout        = "json-rpc-debug.http-timeout"
	FlagJSONRPCDebugHTTPIdleTimeout    = "json-rpc-debug.http-idle-timeout"
	FlagJSONRPCDebugMaxOpenConnections = "json-rpc-debug.max-open-connections"
	FlagJSONRPCDebugReturnDataLimit    = "json-rpc-debug.return-data-limit"
)

const (
	FlagEVMTracer            = "evm.tracer"
	FlagEVMMaxTxGasWanted    = "evm.max-tx-gas-wanted"
	FlagEVMEnableGRPCTracing = "evm.enable-grpc-tracing"
)

var (
	metricsEnabled         bool
	metricsEndpoint        string
	metricsStuckFunc       string
	metricsExportInterval  string
	tracingEnabled         bool
	metricsInsecure        bool
	traceRecorderThreshold int
)

func AddStatsdFlagsToCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&metricsEnabled, "metrics-enable-metrics", false, "Enable OpenTelemetry metrics")
	cmd.PersistentFlags().BoolVar(&tracingEnabled, "metrics-enable-tracing", false, "Enable OpenTelemetry tracing")
	cmd.PersistentFlags().StringVar(&metricsEndpoint, "metrics-endpoint", "localhost:4317", "OpenTelemetry collector gRPC address")
	cmd.PersistentFlags().BoolVar(&metricsInsecure, "metrics-insecure", false, "Disables TLS encryption for gRPC metrics endpoint communication")
	cmd.PersistentFlags().StringVar(&metricsStuckFunc, "metrics-stuck-func", "0m",
		"Sets a duration to consider a function to be stuck to mark in metrics (e.g. in deadlock). 0 disables timeouts.")
	cmd.PersistentFlags().StringVar(&metricsExportInterval, "metrics-export-interval", "10s",
		"Interval to batch and send metrics")
	cmd.PersistentFlags().IntVar(&traceRecorderThreshold, "trace-flight-recorder-threshold", 0,
		"Set trace flight recorder threshold duration in seconds. 0 = flight recorder disabled")
}

func duration(s string, defaults time.Duration) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		dur = defaults
	}
	return dur
}
