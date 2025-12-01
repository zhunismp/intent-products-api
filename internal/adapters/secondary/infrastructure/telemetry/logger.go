package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	otelLog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	dev  = "development"
	prod = "production"
)

func NewLogger(env string) (*zap.Logger, func(context.Context)) {
	logger := zap.Must(zap.NewProduction())
	return logger, func (ctx context.Context) { logger.Sync() }
}

func newProductionLogger(ctx context.Context, appCfg *config.AppEnvConfig) (*otelzap.Logger, func(), error) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   appCfg.GetLogFilePath(),
		MaxSize:    appCfg.GetMaxSize(),
		MaxBackups: appCfg.GetMaxBackups(),
		MaxAge:     appCfg.GetMaxAge(),
		Compress:   appCfg.GetCompress(),
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(lumberjackLogger), zapcore.InfoLevel)
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	core := zapcore.NewTee(fileCore, consoleCore)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(appCfg.GetServerName()),
			semconv.DeploymentEnvironment(appCfg.GetServerEnv()),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(appCfg.GetLogEndpoint()),
		otlploghttp.WithURLPath(appCfg.GetLogPath()),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	processor := otelLog.NewBatchProcessor(exporter)
	provider := otelLog.NewLoggerProvider(
		otelLog.WithResource(res),
		otelLog.WithProcessor(processor),
	)

	global.SetLoggerProvider(provider)

	otelLogger := otelzap.New(zapLogger,
		otelzap.WithLoggerProvider(provider),
		otelzap.WithMinLevel(zapcore.InfoLevel),
	)

	cleanup := func() {
		_ = otelLogger.Sync()
		_ = lumberjackLogger.Close()
		_ = provider.Shutdown(context.Background())
	}

	return otelLogger, cleanup, nil
}

func newDevelopmentLogger() (*otelzap.Logger, func()) {
	otelLogger := otelzap.New(zap.Must(zap.NewDevelopment()))
	return otelLogger, func() { otelLogger.Sync() }
}

func NewLoggerFactory(ctx context.Context, appCfg *config.AppEnvConfig) (*otelzap.Logger, func()) {
	switch appCfg.GetServerEnv() {
	case dev:
		return newDevelopmentLogger()
	case prod:
		logger, cleanupFn, err := newProductionLogger(ctx, appCfg)
		if err != nil {
			log.Fatalf("ERROR: error starting otel logging - %v", err)
		}
		return logger, cleanupFn
	default:
		log.Fatal("error mismatch env type")
		return nil, func() {}
	}

}
