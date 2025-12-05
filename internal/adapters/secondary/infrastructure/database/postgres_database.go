package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/cause"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/product"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// TODO: update gorm to v2
func NewPostgresDatabase(
	host, user, password, dbname, port, sslmode, timezone string,
) (*gorm.DB, func(ctx context.Context) error, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	slog.Info("connecting to postgresql")

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)

	gormConfig := &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, func(ctx context.Context) error { return nil }, fmt.Errorf("open DB: %w", err)
	}

	if err := gormDB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return nil, func(ctx context.Context) error { return nil }, fmt.Errorf("otel plugin: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, func(ctx context.Context) error { return nil }, fmt.Errorf("sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, func(ctx context.Context) error { return sqlDB.Close() }, fmt.Errorf("ping: %w", err)
	}

	slog.Info("database connection established.")

	if err := gormDB.AutoMigrate(
		&ProductModel{},
		&CauseModel{},
	); err != nil {
		return nil, func(ctx context.Context) error { return sqlDB.Close() }, fmt.Errorf("auto-migrate: %w", err)
	}

	shutdownFn := func(ctx context.Context) error {
		done := make(chan error, 1)

		go func() {
			done <- sqlDB.Close()
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-done:
			return err
		}
	}

	return gormDB, shutdownFn, nil
}
