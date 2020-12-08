package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

// NewLogger used to initialize the application logger
func NewLogger(debugMode bool) *zap.SugaredLogger {
	var level zapcore.Level

	switch debugMode {
	case true:
		level = zapcore.DebugLevel
	default:
		level = zapcore.InfoLevel
	}

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			LevelKey:    "key",
		},
	}

	logger, _ := cfg.Build()
	appLogger = logger.Sugar()

	return appLogger
}

// Get used to get the application logger
func Get() *zap.SugaredLogger {
	defer appLogger.Sync() // nolint: errcheck
	return appLogger
}

// ExitWithError will terminate execution with an error result
// It prints the error to stderr and exits with a non-zero exit code
func ExitWithError(err error) {
	defer appLogger.Sync() // nolint: errcheck
	appLogger.Error(err)
	os.Exit(1)
}
