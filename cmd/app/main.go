package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http"
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/config"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/database"
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

	db := NewPostgresDatabase(
		cfg.GetDBHost(),
		cfg.GetDBUser(),
		cfg.GetDBPassword(),
		cfg.GetDBName(),
		cfg.GetDBPort(),
		cfg.GetDBSSLMode(),
		cfg.GetDBTimezone(),
	)
	_, err = SetupTelemetry(context.Background(), cfg.GetServerName(), cfg.GetServerEnv())
	
	logger := GetLogger(cfg.GetServerEnv(), cfg.GetServerEnv())
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	logger.Info(fmt.Sprintf("Received shutdown signal %s", sig.String()))
	httpServer.GracefulShutdown()
	logger.Info("Cleanup finished. Exiting...")
}
