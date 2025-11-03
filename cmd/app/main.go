package main

import (
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http"
	. "github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/infrastructure/database"
	. "github.com/zhunismp/intent-products-api/internal/adapters/secondary/repositories/cause"
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

	db := NewPostgresDatabase(
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

	productDbRepo := NewProductRepository(db)
	causeDbRepo := NewCauseRepository(db)
	productSvc := NewProductService(productDbRepo, causeDbRepo)
	productHttp := NewProductHttpHandler(productSvc)

	routeGroup := NewRouteGroup(productHttp)

	httpServer := NewHttpServer(cfg, log, baseApiPrefix)
	httpServer.SetupRoute(routeGroup)
	httpServer.Start()
}
