package logger

import (
	"log"

	"go.uber.org/zap"
)

const (
	dev  = "development"
	prod = "production"
)

func NewLogger(env string) *zap.Logger {
	return newLoggerFactory(env)
}

func newLoggerFactory(env string) *zap.Logger {
	switch env {
	case dev:
		return zap.Must(zap.NewDevelopment())
	case prod:
		return zap.Must(zap.NewProduction())
	default:
		log.Fatal("error mismatch env type")
		return nil
	}
}
