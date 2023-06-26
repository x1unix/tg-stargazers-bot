package config

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Environment Environment   `envconfig:"APP_ENV" default:"dev"`
	Level       zapcore.Level `envconfig:"LOG_LEVEL" default:"debug"`
	Format      string        `envconfig:"LOG_FORMAT" default:"text"`
}

func (cfg LogConfig) NewLogger() (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	logCfg.Development = cfg.Environment != ProdEnvironment
	logCfg.Level = zap.NewAtomicLevelAt(cfg.Level)
	logCfg.Encoding = cfg.Format

	switch cfg.Format {
	case "", "json":
		logCfg.EncoderConfig = zap.NewProductionEncoderConfig()
	case "text":
		logCfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	default:
		return nil, fmt.Errorf("unsupported log format %q", cfg.Format)
	}

	return logCfg.Build()
}
