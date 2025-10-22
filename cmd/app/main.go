package main

import (
	"context"
	"time"

	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http"
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/database"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/product"
	. "github.com/zhunismp/intent-products-api/internal/core/domain/product"
	. "github.com/zhunismp/intent-products-api/internal/infrastructure/config"
	. "github.com/zhunismp/intent-products-api/internal/infrastructure/logger"
)

func main() {
	cfg, err := LoadConfig(".env")
	if err != nil {
		panic(err.Error())
	}

	initDbCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	db, closeConn := NewMongoDatabase(
		initDbCtx,
		cfg.GetDBHost(),
		cfg.GetDBUser(),
		cfg.GetDBPassword(),
		cfg.GetDBName(),
		cfg.GetDBPort(),
		cfg.GetDBSSLMode(),
		cfg.GetDBTimezone(),
	)
	log := NewLogger(cfg.GetServerEnv())
	baseApiPrefix := cfg.GetServerBaseApiPrefix()

	productDbRepo := NewProductRepository(initDbCtx, db)
	productSvc := NewProductService(productDbRepo)
	productHttp := NewProductHttpHandler(productSvc)

	routeGroup := NewRouteGroup(productHttp)

	httpServer := NewHttpServer(cfg, log, baseApiPrefix)
	httpServer.SetupRoute(routeGroup)
	httpServer.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		// Other close connection func
		log.Info("shutting down external dependencies connection.")
		closeConn(ctx)
		cancel()
	}()
}
