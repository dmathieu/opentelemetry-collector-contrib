// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package datastream // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter/internal/datastream"

const (
	DatasetKey   = "data_stream.dataset"
	NamespaceKey = "data_stream.namespace"
	TypeKey      = "data_stream.type"

	defaultDataset     = "generic"
	defaultNamespace   = "default"
	defaultTypeLogs    = "logs"
	defaultTypeMetrics = "metrics"
	defaultTypeTraces  = "traces"
)
