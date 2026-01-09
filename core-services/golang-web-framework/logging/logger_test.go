package logging

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/stretchr/testify/assert"
)

func TestInfoLevelLog(t *testing.T) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")
	configManager := config.GetConfigurationManagerInstance()
	logFileLocation := configManager.GetConfigFor("logger.file-name")

	tearDown(logFileLocation)

	var logger = GetLoggerInstance()

	logger.LogInfoFor("A Message")

	content, err := os.ReadFile(logFileLocation)
	if err != nil {
		t.Fatal("no log file was written", err)
	}

	actual := LogRecord{}

	err = json.Unmarshal(content, &actual)
	if err != nil {
		t.Fatal("log record parsing error: ", err)

	}
	assert.Equal(t, "A Message", actual.Msg)
	assert.Equal(t, "info", actual.Level)

}

func tearDown(logFileLocation string) {
	_, error := os.Stat(logFileLocation)
	if error != nil {
		os.Remove(logFileLocation)
	}
}

type LogRecord struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}
