package logging

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"gopkg.in/natefinch/lumberjack.v2"

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

func GetLoggerInstanceForComponentByTypeName(className string) *Logger {
	manager = config.GetConfigurationManagerInstance()
	return &Logger{
		logger: new().logger.With("component", className),
	}
}

func GetLoggerInstanceForComponentByType(obj any) *Logger {
	manager = config.GetConfigurationManagerInstance()
	typeName := reflect.TypeOf(obj).Elem().Name()

	return &Logger{
		logger: new().logger.With("component", typeName),
	}
}

func resolveLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func logInit(f *os.File) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := resolveLevel(manager.GetConfigFor("logger.level"))

	var cores []zapcore.Core
	if f == nil {
		cores = []zapcore.Core{
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		}
	} else {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   f.Name(),
			MaxSize:    1,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		}
	}

	if manager.GetConfigBoolFor("otel.enabled") {
		serviceName := manager.GetConfigFor("otel.service-name")
		cores = append(cores, otelzap.NewCore(serviceName))
	}

	l := zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return l.Sugar()
}

func new() *Logger {
	fileName := manager.GetConfigFor("logger.file-name")
	var f *os.File
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("log file does not exist")
	}

	return &Logger{
		logger: logInit(f),
	}
}

func (impl *Logger) LogErrorFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Error(str)
}

func (impl *Logger) LogErrorfFor(format string, a ...any) {
	impl.logger.Errorf(format, a...)
}

func (impl *Logger) LogInfoFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Info(str)
}

func (impl *Logger) LogInfofFor(format string, a ...any) {
	impl.logger.Infof(format, a...)
}

func (impl *Logger) LogDebugFor(message any) {
	str := fmt.Sprintf("%v", message)
	impl.logger.Debug(str)
}

func (impl *Logger) LogDebugfFor(format string, a ...any) {
	impl.logger.Debugf(format, a...)
}

// todo add a function to flush the logger buffer, to be called on application shutdown
func Dispose() {
	if loggerManager != nil {
		loggerManager.logger.Sync()
	}
}
