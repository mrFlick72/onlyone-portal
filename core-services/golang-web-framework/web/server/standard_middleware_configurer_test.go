package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// k8s liveness/readiness probes hit healthCheckPath every few seconds per
// pod; the default must skip it from the gin access log so probe traffic
// doesn't drown out real requests (see ADR 0005).
func TestAccessLogSkipPaths_DefaultSkipsHealthCheck(t *testing.T) {
	assert.Equal(t, []string{healthCheckPath}, accessLogSkipPaths(false))
}

// server.access-log.health-check-logging-enabled: true is the opt-in escape
// hatch to see health-check requests in the access log again.
func TestAccessLogSkipPaths_EnabledLogsHealthCheck(t *testing.T) {
	assert.Nil(t, accessLogSkipPaths(true))
}
