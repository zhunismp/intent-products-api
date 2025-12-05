package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http"
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/config"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/database"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/shutdown"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/telemetry"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/cause"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/product"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/product"
)

func main() {
	cfg, err := LoadConfig(".env")
	if err != nil {
		panic(err.Error())
	}
	logger := GetLogger(cfg.GetServerEnv(), cfg.GetServerName())
	sm := NewShutdownManager(20*time.Second, logger)

	otelShutdownFn, err := SetupTelemetry(context.Background(), cfg.GetServerName(), cfg.GetServerEnv())
	if err != nil {
		sm.Shutdown()
	}
	sm.Register(&ShutdownFunction{
		ResourceName: "opentelemetry",
		Fn:           otelShutdownFn,
	})

	db, dbShutdownFn, err := NewPostgresDatabase(
		cfg.GetDBHost(),
		cfg.GetDBUser(),
		cfg.GetDBPassword(),
		cfg.GetDBName(),
		cfg.GetDBPort(),
		cfg.GetDBSSLMode(),
		cfg.GetDBTimezone(),
	)
	if err != nil {
		sm.Shutdown()
	}
	sm.Register(&ShutdownFunction{
		ResourceName: "database",
		Fn:           dbShutdownFn,
	})

	baseApiPrefix := cfg.GetServerBaseApiPrefix()

	productDbRepo := NewProductRepository(db)
	causeDbRepo := NewCauseRepository(db)

	causeSvc := NewCauseService(causeDbRepo, logger)
	productSvc := NewProductService(productDbRepo, causeSvc, logger)

	// HTTP
	productHttp := NewProductHttpHandler(productSvc, logger)
	routeGroup := NewRouteGroup(productHttp)
	httpServer := NewHttpServer(cfg, logger, baseApiPrefix)
	httpServer.SetupRoute(routeGroup)
	httpServer.Start()
	sm.Register(&ShutdownFunction{
		ResourceName: "http server",
		Fn:           httpServer.GracefulShutdown,
	})

	gracefulShutdown(sm)
}

func gracefulShutdown(sm ShutdownManager) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit

	sm.Shutdown()
}
