// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package elasticsearchexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterbatcher"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/internal/config"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/internal/metadata"
)

// NewFactory creates a factory for Elastic exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		config.CreateDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
	)
}

// createLogsExporter creates a new exporter for logs.
//
// Logs are directly indexed into Elasticsearch.
func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	cf := cfg.(*config.Config)

	index := cf.LogsIndex
	if cf.Index != "" {
		set.Logger.Warn("index option are deprecated and replaced with logs_index and traces_index.")
		index = cf.Index
	}
	config.HandleDeprecatedConfig(cf, set.Logger)

	exporter := newExporter(cf, set, index, cf.LogsDynamicIndex.Enabled)

	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		exporter.pushLogsData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	cf := cfg.(*config.Config)
	config.HandleDeprecatedConfig(cf, set.Logger)

	exporter := newExporter(cf, set, cf.MetricsIndex, cf.MetricsDynamicIndex.Enabled)

	return exporterhelper.NewMetrics(
		ctx,
		set,
		cfg,
		exporter.pushMetricsData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func createTracesExporter(ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Traces, error) {
	cf := cfg.(*config.Config)
	config.HandleDeprecatedConfig(cf, set.Logger)

	exporter := newExporter(cf, set, cf.TracesIndex, cf.TracesDynamicIndex.Enabled)

	return exporterhelper.NewTraces(
		ctx,
		set,
		cfg,
		exporter.pushTraceData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func exporterhelperOptions(
	cfg *config.Config,
	start component.StartFunc,
	shutdown component.ShutdownFunc,
) []exporterhelper.Option {
	opts := []exporterhelper.Option{
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
		exporterhelper.WithStart(start),
		exporterhelper.WithShutdown(shutdown),
		exporterhelper.WithQueue(cfg.QueueSettings),
	}
	if cfg.Batcher.Enabled != nil {
		batcherConfig := exporterbatcher.Config{
			Enabled:       *cfg.Batcher.Enabled,
			FlushTimeout:  cfg.Batcher.FlushTimeout,
			MinSizeConfig: cfg.Batcher.MinSizeConfig,
			MaxSizeConfig: cfg.Batcher.MaxSizeConfig,
		}
		opts = append(opts, exporterhelper.WithBatcher(batcherConfig))

		// Effectively disable timeout_sender because timeout is enforced in bulk indexer.
		//
		// We keep timeout_sender enabled in the async mode (Batcher.Enabled == nil),
		// to ensure sending data to the background workers will not block indefinitely.
		opts = append(opts, exporterhelper.WithTimeout(exporterhelper.TimeoutConfig{Timeout: 0}))
	}
	return opts
}
