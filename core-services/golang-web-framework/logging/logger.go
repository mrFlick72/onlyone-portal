package logging

import (
	"fmt"
	"os"
	"sync"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	manager       *config.ConfigurationManager
	loggerManager *Logger
	once          sync.Once
)

type Logger struct {
	logger *zap.SugaredLogger
}

func GetLoggerInstance() *Logger {
	once.Do(func() {
		manager = config.GetConfigurationManagerInstance()

		loggerManager = new()
	})

	return loggerManager
}
func logInit(f *os.File) *zap.SugaredLogger {

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core)

	return l.Sugar()
}

func new() *Logger {
	fileName := manager.GetConfigFor("logger.file-name")
	var f *os.File
	f, err := os.Create(fileName)
	if err != nil {
		panic("log file does not exist")
	}

	return &Logger{
		logger: logInit(f),
	}
}

func (impl *Logger) LogErrorFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Error(str)
}

func (impl *Logger) LogInfoFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Info(str)
}
func (impl *Logger) LogDebugFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Debug(str)
}

func Dispose() {
	panic("TODO it have to be implemented")
}
