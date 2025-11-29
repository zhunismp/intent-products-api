package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http"
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/database"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/cause"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/product"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/cause"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/priority"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	. "github.com/zhunismp/intent-products-api/internal/infrastructure/config"
	. "github.com/zhunismp/intent-products-api/internal/infrastructure/logger"
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
	logger := NewLogger(cfg.GetServerEnv())
	baseApiPrefix := cfg.GetServerBaseApiPrefix()

	productDbRepo := NewProductRepository(db, logger)
	causeDbRepo := NewCauseRepository(db, logger)

	causeSvc := NewCauseService(causeDbRepo, logger)
	prioritySvc := NewPriorityService(logger)
	productSvc := NewProductService(productDbRepo, causeSvc, prioritySvc, logger)

	// HTTP
	productHttp := NewProductHttpHandler(productSvc)
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
