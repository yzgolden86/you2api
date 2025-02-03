package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) error {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(getLogLevel(level))
	
	logger, err := config.Build()
	if err != nil {
		return err
	}
	
	log = logger
	return nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	// ... 其他级别
	default:
		return zapcore.InfoLevel
	}
} 