// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/config"

import (
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterbatcher"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	defaultLogsIndex    = "logs-generic-default"
	defaultMetricsIndex = "metrics-generic-default"
	defaultTracesIndex  = "traces-generic-default"
)

func CreateDefaultConfig() component.Config {
	qs := exporterhelper.NewDefaultQueueConfig()
	qs.Enabled = false

	httpClientConfig := confighttp.NewDefaultClientConfig()
	httpClientConfig.Timeout = 90 * time.Second
	httpClientConfig.Compression = configcompression.TypeGzip

	return &Config{
		QueueSettings: qs,
		ClientConfig:  httpClientConfig,
		Index:         "",
		LogsIndex:     defaultLogsIndex,
		LogsDynamicIndex: DynamicIndexSetting{
			Enabled: false,
		},
		MetricsIndex: defaultMetricsIndex,
		MetricsDynamicIndex: DynamicIndexSetting{
			Enabled: true,
		},
		TracesIndex: defaultTracesIndex,
		TracesDynamicIndex: DynamicIndexSetting{
			Enabled: false,
		},
		Retry: RetrySettings{
			Enabled:         true,
			MaxRetries:      0, // default is set in exporter code
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     1 * time.Minute,
			RetryOnStatus: []int{
				http.StatusTooManyRequests,
			},
		},
		Mapping: MappingsSettings{
			Mode:  "none",
			Dedot: true,
		},
		LogstashFormat: LogstashFormatSettings{
			Enabled:         false,
			PrefixSeparator: "-",
			DateFormat:      "%Y.%m.%d",
		},
		TelemetrySettings: TelemetrySettings{
			LogRequestBody:  false,
			LogResponseBody: false,
		},
		Batcher: BatcherConfig{
			FlushTimeout: 30 * time.Second,
			MinSizeConfig: exporterbatcher.MinSizeConfig{
				MinSizeItems: 5000,
			},
			MaxSizeConfig: exporterbatcher.MaxSizeConfig{
				MaxSizeItems: 0,
			},
		},
		Flush: FlushSettings{
			Bytes:    5e+6,
			Interval: 30 * time.Second,
		},
	}
}

func WithDefaultConfig(fns ...func(*Config)) *Config {
	cfg := CreateDefaultConfig().(*Config)
	for _, fn := range fns {
		fn(cfg)
	}
	return cfg
}

func WithDefaultHTTPClientConfig(fns ...func(config *confighttp.ClientConfig)) confighttp.ClientConfig {
	cfg := confighttp.NewDefaultClientConfig()
	for _, fn := range fns {
		fn(&cfg)
	}
	return cfg
}
