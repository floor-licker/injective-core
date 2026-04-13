package main

import (
	"time"

	"github.com/InjectiveLabs/metrics/v2"
	"github.com/xlab/closer"
)

var (
	metricsEnabled        *bool
	metricsEndpoint       *string
	metricsInsecure       *bool
	metricsExportInterval *string
	metricsStuckFunc      *string
	tracingEnabled        *bool
)

func initMetrics(chainID string) (metrics.Meter, error) {
	stuckFuncTimeout, err := time.ParseDuration(*metricsStuckFunc)
	if err != nil {
		return nil, err
	}

	exportInterval, err := time.ParseDuration(*metricsExportInterval)
	if err != nil {
		return nil, err
	}

	cfg := metrics.Config{
		Endpoint:         *metricsEndpoint,
		InsecureEndpoint: *metricsInsecure,
		MetricsEnabled:   *metricsEnabled,
		TracingEnabled:   *tracingEnabled,
		StuckFuncTimeout: stuckFuncTimeout,
		ExportInterval:   exportInterval,
	}

	resourceAttrs := []metrics.TagAttribute{
		metrics.Tag(metrics.ServiceNameKey, "peggo"),
	}
	if chainID != "" {
		resourceAttrs = append(resourceAttrs, metrics.Tag("chain-id", chainID))
	}
	if envName != nil && *envName != "" {
		resourceAttrs = append(resourceAttrs, metrics.Tag("env", *envName))
	}

	peggoMetrics, err := metrics.NewMetrics(cfg, resourceAttrs...)
	if err != nil {
		return nil, err
	}

	peggoMeter, err := peggoMetrics.NewMeter("peggo")
	if err != nil {
		return nil, err
	}

	closer.Bind(func() {
		_ = peggoMetrics.Shutdown()
	})

	return peggoMeter, nil
}
