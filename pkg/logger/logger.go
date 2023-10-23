package logger

import (
	"go.uber.org/zap"
)

func InitLogger(level string) *zap.Logger {
	logLevel := zap.ErrorLevel

	if level == "DEBUG" {
		logLevel = zap.DebugLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)

	return logger
}
