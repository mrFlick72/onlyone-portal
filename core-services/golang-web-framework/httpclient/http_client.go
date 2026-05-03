package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

// NewHTTPClient returns an *http.Client. When otel.enabled=true the transport
// injects W3C traceparent headers and creates a client span for every outgoing
// request. The request must be created with http.NewRequestWithContext so the
// active span is reachable from the transport. Optional otelhttp options let
// callers customize span names or attributes for a specific outbound client.
// When otel.enabled=false a plain http.Client is returned with no overhead.
func NewHTTPClient(options ...otelhttp.Option) *http.Client {
	if config.GetConfigurationManagerInstance().GetConfigBoolFor("otel.enabled") {
		return &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport, options...),
		}
	}
	return &http.Client{}
}
