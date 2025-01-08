// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/config"

import "go.uber.org/zap"

func HandleDeprecatedConfig(cfg *Config, logger *zap.Logger) {
	if cfg.Mapping.Dedup != nil {
		logger.Warn("dedup is deprecated, and is always enabled")
	}
	if cfg.Mapping.Dedot && cfg.MappingMode() != MappingECS || !cfg.Mapping.Dedot && cfg.MappingMode() == MappingECS {
		logger.Warn("dedot has been deprecated: in the future, dedotting will always be performed in ECS mode only")
	}
	if cfg.Retry.MaxRequests != 0 {
		cfg.Retry.MaxRetries = cfg.Retry.MaxRequests - 1
		// Do not set cfg.Retry.Enabled = false if cfg.Retry.MaxRequest = 1 to avoid breaking change on behavior
		logger.Warn("retry::max_requests has been deprecated, and will be removed in a future version. Use retry::max_retries instead.")
	}
}
