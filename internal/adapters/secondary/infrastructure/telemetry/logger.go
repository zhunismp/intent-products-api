package telemetry

import (
	"log"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	dev  = "development"
	prod = "production"
)

func NewLoggerFactory(env string) (*zap.Logger, func() error) {
	switch env {
	case prod:
		logger := otelzap.New(zap.Must(zap.NewProduction()))
		return logger.Logger, logger.Sync
	case dev:
		logger := zap.Must(zap.NewDevelopment())
		return logger, logger.Sync
	default:
		log.Fatalf("environment configuration mismatch")
		return nil, nil
	}
}
