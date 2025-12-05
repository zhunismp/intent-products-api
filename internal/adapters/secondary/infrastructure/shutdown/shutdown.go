package shutdown

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"
)

type ShutdownFunction struct {
	ResourceName string
	Fn           func(ctx context.Context) error
}

type ShutdownManager interface {
	Register(fn *ShutdownFunction)
	Shutdown()
}

type shutdownManagerImpl struct {
	timeout     time.Duration
	shutdownFns []*ShutdownFunction
	logger      *slog.Logger
}

func NewShutdownManager(timeout time.Duration, logger *slog.Logger) ShutdownManager {
	return &shutdownManagerImpl{
		timeout:     timeout,
		shutdownFns: make([]*ShutdownFunction, 0),
		logger:      logger,
	}
}

func (s *shutdownManagerImpl) Register(fn *ShutdownFunction) {
	s.shutdownFns = append(s.shutdownFns, fn)
	s.logger.Info("registered shutdown function", slog.String("resource", fn.ResourceName))
}

func (s *shutdownManagerImpl) Shutdown() {
	s.logger.Info("starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var errs []error
	for i := len(s.shutdownFns) - 1; i >= 0; i-- {
		sf := s.shutdownFns[i]
		if err := sf.Fn(ctx); err != nil {
			s.logger.Error(
				"resource cleanup failed",
				slog.String("resource", sf.ResourceName),
				slog.Any("error", err),
			)
			errs = append(errs, err)
		} else {
			s.logger.Info(
				"resource cleaned successfully",
				slog.String("resource", sf.ResourceName),
			)
		}
	}

	if len(errs) > 0 {
		agg := errors.Join(errs...)
		s.logger.Error("shutdown completed with errors", slog.Any("errors", agg))
		os.Exit(1)
	}

	s.logger.Info("shutdown completed successfully")
	os.Exit(0)
}
