package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is a
type Config struct {
	Debug       bool   `envconfig:"debug" default:"0"`
	Tags        string `envconfig:"tags" default:"Name=*"`
	SSHUsername string `envconfig:"ssh_username" default:"ec2-user"`
	SSHPort     string `envconfig:"ssh_port" default:"22"`
	SSHOpts     string `envconfig:"ssh_opts" default:"-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"`
}

var appConfig Config
var appLogger *zap.SugaredLogger

// Load is used to load the application configuration
func Load() {
	if err := envconfig.Process("awssh", &appConfig); err != nil {
		log.Fatal(err.Error())
	}
}

// Get is used to gather the application configuration
func Get() *Config {
	return &appConfig
}

// LoadLogger is
func LoadLogger() *zap.SugaredLogger {
	var level zapcore.Level

	switch appConfig.Debug {
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

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()
	appLogger = logger.Sugar()

	return appLogger
}

// GetLogger is
func GetLogger() *zap.SugaredLogger {
	return appLogger
}
