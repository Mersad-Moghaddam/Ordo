package logging

import (
	"go.uber.org/zap"
)

type LoggerConfiguration struct {
	DevelopmentMode bool
}

type Option func(loggerConfiguration *LoggerConfiguration)

func NewZapLogger(options ...Option) (*zap.Logger, error) {
	loggerConfiguration := LoggerConfiguration{DevelopmentMode: false}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&loggerConfiguration)
	}

	if loggerConfiguration.DevelopmentMode {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func WithDevelopmentMode() Option {
	return func(loggerConfiguration *LoggerConfiguration) {
		loggerConfiguration.DevelopmentMode = true
	}
}
