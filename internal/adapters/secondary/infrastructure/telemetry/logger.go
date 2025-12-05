package telemetry

import (
	"log"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	slogctx "github.com/veqryn/slog-context"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

const (
	development = "development"
	production  = "production"
)

func NewProductionLogger(appName string) *slog.Logger {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	otelHandler := otelslog.NewHandler(appName, otelslog.WithSource(true))
	fanoutHandler := slogmulti.Fanout(stdoutHandler, otelHandler)
	ctxHandler := slogctx.NewHandler(fanoutHandler, nil)
	logger := slog.New(ctxHandler)
	return logger
}

func NewDevelopmentLogger(appName string) *slog.Logger {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	otelHandler := otelslog.NewHandler(appName, otelslog.WithSource(true))
	fanoutHandler := slogmulti.Fanout(stdoutHandler, otelHandler)
	ctxHandler := slogctx.NewHandler(fanoutHandler, nil)
	logger := slog.New(ctxHandler)
	return logger
}

func GetLogger(env, appName string) *slog.Logger {
	switch env {
	case development:
		return NewDevelopmentLogger(appName)
	case production:
		return NewProductionLogger(appName)
	default:
		log.Fatalf("wrong environment was set. application can not proceed")
		return nil
	}
}
