package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger *zap.SugaredLogger

// Init used to initialize the application logger
func Init(debugMode bool) {
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

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, _ := cfg.Build()
	appLogger = logger.Sugar()
}

// Get used to get the application logger
func Get() *zap.SugaredLogger {
	defer appLogger.Sync()
	return appLogger
}
