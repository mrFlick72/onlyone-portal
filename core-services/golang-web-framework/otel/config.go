package otel

import "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"

type otelConfig struct {
	Enabled     bool
	ServiceName string
	Protocol    string // "http" (default) or "grpc"
	Endpoint    string // OTel Collector host:port, used for all signals
	Insecure    bool
}

func loadOtelConfig() otelConfig {
	mgr := config.GetConfigurationManagerInstance()
	return otelConfig{
		Enabled:     mgr.GetConfigBoolFor("otel.enabled"),
		ServiceName: mgr.GetConfigFor("otel.service-name"),
		Protocol:    mgr.GetConfigFor("otel.protocol"),
		Endpoint:    mgr.GetConfigFor("otel.endpoint"),
		Insecure:    mgr.GetConfigBoolFor("otel.insecure"),
	}
}
