package otel

import "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"

type signalConfig struct {
	Endpoint string
}

type otelConfig struct {
	Enabled     bool
	ServiceName string
	Protocol    string // "http" (default) or "grpc"
	Insecure    bool
	Traces      signalConfig
	Metrics     signalConfig
	Logs        signalConfig
}

func loadOtelConfig() otelConfig {
	mgr := config.GetConfigurationManagerInstance()
	return otelConfig{
		Enabled:     mgr.GetConfigBoolFor("otel.enabled"),
		ServiceName: mgr.GetConfigFor("otel.service-name"),
		Protocol:    mgr.GetConfigFor("otel.protocol"),
		Insecure:    mgr.GetConfigBoolFor("otel.insecure"),
		Traces: signalConfig{
			Endpoint: mgr.GetConfigFor("otel.traces.endpoint"),
		},
		Metrics: signalConfig{
			Endpoint: mgr.GetConfigFor("otel.metrics.endpoint"),
		},
		Logs: signalConfig{
			Endpoint: mgr.GetConfigFor("otel.logs.endpoint"),
		},
	}
}
