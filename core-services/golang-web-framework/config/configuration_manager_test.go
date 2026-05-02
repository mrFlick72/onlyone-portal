package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")

	var configManager = GetConfigurationManagerInstance()

	var serverPortConfigurationValue = configManager.GetConfigFor("server.port")

	assert.Equal(t, "3050", serverPortConfigurationValue)
}

func TestGetConfigDurationFor_ReturnsConfiguredValue(t *testing.T) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")
	configManager := GetConfigurationManagerInstance()

	got := configManager.GetConfigDurationFor("server.read-timeout", 99*time.Second)

	assert.Equal(t, 45*time.Second, got)
}

func TestGetConfigDurationFor_ReturnsDefaultWhenKeyMissing(t *testing.T) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")
	configManager := GetConfigurationManagerInstance()

	got := configManager.GetConfigDurationFor("server.no-such-key", 7*time.Second)

	assert.Equal(t, 7*time.Second, got)
}
