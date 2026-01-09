package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestHealth(t *testing.T) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")

	var configManager = GetConfigurationManagerInstance()

	var serverPortConfigurationValue = configManager.GetConfigFor("server.port")

	assert.Equal(t, "3050", serverPortConfigurationValue)
}
