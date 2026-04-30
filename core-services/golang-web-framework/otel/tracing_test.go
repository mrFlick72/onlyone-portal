package otel

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func TestMain(m *testing.M) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")
	os.Exit(m.Run())
}

// --- Setup (unified) ---

func TestSetup_WhenDisabled_ReturnsNoopAndInstallsNothing(t *testing.T) {
	originalTP := otel.GetTracerProvider()
	originalMP := otel.GetMeterProvider()
	t.Cleanup(func() {
		otel.SetTracerProvider(originalTP)
		otel.SetMeterProvider(originalMP)
	})

	shutdown, err := Setup(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	assert.Equal(t, originalTP, otel.GetTracerProvider())
	assert.Equal(t, originalMP, otel.GetMeterProvider())
	assert.NoError(t, shutdown(context.Background()))
}

// --- TracerProvider ---

func TestSetupTracerProvider_WhenDisabled_ReturnsNoop(t *testing.T) {
	original := otel.GetTracerProvider()
	t.Cleanup(func() { otel.SetTracerProvider(original) })

	shutdown, err := SetupTracerProvider(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	assert.Equal(t, original, otel.GetTracerProvider())
	assert.NoError(t, shutdown(context.Background()))
}

func TestBuildExporter_HTTP_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "http", Traces: signalConfig{Endpoint: "localhost:4318"}, Insecure: true}
	exp, err := buildExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildExporter_GRPC_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "grpc", Traces: signalConfig{Endpoint: "localhost:4317"}, Insecure: true}
	exp, err := buildExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildExporter_DefaultsToHTTP(t *testing.T) {
	cfg := otelConfig{Protocol: "", Traces: signalConfig{Endpoint: "localhost:4318"}, Insecure: true}
	exp, err := buildExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildExporter_UnsupportedProtocol_ReturnsError(t *testing.T) {
	cfg := otelConfig{Protocol: "kafka", Traces: signalConfig{Endpoint: "localhost:9092"}}
	exp, err := buildExporter(context.Background(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol")
	assert.Nil(t, exp)
}

// --- MeterProvider ---

func TestBuildMetricExporter_HTTP_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "http", Metrics: signalConfig{Endpoint: "localhost:4318"}, Insecure: true}
	exp, err := buildMetricExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildMetricExporter_GRPC_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "grpc", Metrics: signalConfig{Endpoint: "localhost:4317"}, Insecure: true}
	exp, err := buildMetricExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildMetricExporter_UnsupportedProtocol_ReturnsError(t *testing.T) {
	cfg := otelConfig{Protocol: "kafka", Metrics: signalConfig{Endpoint: "localhost:9092"}}
	exp, err := buildMetricExporter(context.Background(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol")
	assert.Nil(t, exp)
}

// --- LoggerProvider ---

func TestBuildLogExporter_HTTP_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "http", Logs: signalConfig{Endpoint: "localhost:4318"}, Insecure: true}
	exp, err := buildLogExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildLogExporter_GRPC_Succeeds(t *testing.T) {
	cfg := otelConfig{Protocol: "grpc", Logs: signalConfig{Endpoint: "localhost:4317"}, Insecure: true}
	exp, err := buildLogExporter(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	_ = exp.Shutdown(context.Background())
}

func TestBuildLogExporter_UnsupportedProtocol_ReturnsError(t *testing.T) {
	cfg := otelConfig{Protocol: "kafka", Logs: signalConfig{Endpoint: "localhost:9092"}}
	exp, err := buildLogExporter(context.Background(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol")
	assert.Nil(t, exp)
}
