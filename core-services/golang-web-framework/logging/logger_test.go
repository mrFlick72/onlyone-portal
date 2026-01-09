package logging

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/stretchr/testify/assert"
)

var logFileLocation string
var uut *Logger

func TestMain(m *testing.M) {
	os.Setenv("CONFIG_FILE_LOCATION", "../test/application.yml")
	configManager := config.GetConfigurationManagerInstance()
	logFileLocation = configManager.GetConfigFor("logger.file-name")

	uut = GetLoggerInstance()

	code := m.Run() // run tests

	// global teardown
	os.Remove(logFileLocation)

	os.Exit(code)
}

func TestInfoLevelLog(t *testing.T) {
	uut.LogInfoFor("A Message")

	actual := getLogRecordFrom(t)

	assert.Equal(t, "A Message", actual.Msg)
	assert.Equal(t, "info", actual.Level)
}

func getLogRecordFrom(t *testing.T) LogRecord {
	content, err := os.ReadFile(logFileLocation)
	if err != nil {
		t.Fatal("no log file was written ", logFileLocation, err)
	}

	actual := LogRecord{}

	err = json.Unmarshal(content, &actual)
	if err != nil {
		t.Fatal("log record parsing error: ", err)

	}

	return actual
}

type LogRecord struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}
