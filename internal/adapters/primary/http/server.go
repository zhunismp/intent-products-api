package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	cors "github.com/gofiber/fiber/v3/middleware/cors"
	limiter "github.com/gofiber/fiber/v3/middleware/limiter"
	recover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/zhunismp/intent-products-api/internal/adapters/primary/http/middleware"
	"github.com/zhunismp/intent-products-api/internal/adapters/primary/http/product"
	core "github.com/zhunismp/intent-products-api/internal/core/infrastructure/config"
)

type HttpServer struct {
	cfg           core.AppConfigProvider
	log           *slog.Logger
	fiberApp      *fiber.App
	apiBaseRouter fiber.Router
	basePath      string
}

type RouteGroup struct {
	product *product.ProductHttpHandler
}

func NewRouteGroup(product *product.ProductHttpHandler) *RouteGroup {
	return &RouteGroup{product: product}
}

func NewHttpServer(cfg core.AppConfigProvider, slogLogger *slog.Logger, baseApiPrefix string) *HttpServer {
	validateArguments(cfg, slogLogger, &baseApiPrefix)

	app := fiber.New(fiber.Config{
		AppName: cfg.GetServerName(),
	})

	app.Use(recover.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.TraceMiddleware())
	app.Use(middleware.AccessLogMiddleware(slogLogger))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"errorMessage": "too many requests, please try again later.",
			})
		},
	}))

	apiGroup := app.Group(baseApiPrefix)
	slogLogger.Info("fiber HTTP server core initialized with middleware.", slog.String("baseApiPrefix", baseApiPrefix))

	return &HttpServer{
		cfg:           cfg,
		log:           slogLogger,
		fiberApp:      app,
		apiBaseRouter: apiGroup,
		basePath:      baseApiPrefix,
	}
}

func (s *HttpServer) Start() {
	serverAddr := fmt.Sprintf("%s:%s", s.cfg.GetServerHost(), s.cfg.GetServerPort())
	s.log.Info("attempting to start HTTP server...", slog.String("address", serverAddr))

	go func() {
		if err := s.fiberApp.Listen(serverAddr); err != nil && err.Error() != "http: Server closed" {
			s.log.Error("failed to start HTTP server listener")
		}
	}()

	s.log.Info("HTTP server listener started.", slog.String("address", serverAddr))
}

func (s *HttpServer) SetupRoute(routeGroup *RouteGroup) {
	if routeGroup.product == nil {
		s.log.Error("failed to set up route")
	}

	productHandler := routeGroup.product

	s.registerAPIGroup("/health", func(router fiber.Router) {
		router.Get("/liveness", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"message": "application is running"}) })
	})

	s.registerAPIGroup("/products", func(router fiber.Router) {
		// core product
		router.Get("/:id", productHandler.GetProduct)
		router.Get("/", productHandler.GetAllProducts)
		router.Post("/", productHandler.CreateProduct)
		router.Put("/positions", productHandler.MoveProductPosition)
		router.Delete("/:id", productHandler.DeleteProduct)
	})
}

func (s *HttpServer) GracefulShutdown(ctx context.Context) error {
	s.log.Info("gracefully shutting down HTTP server...")

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		s.log.Error("error during server shutdown")
		return err
	} else {
		s.log.Info("HTTP server shutdown gracefully.")
		return nil
	}
}

func (s *HttpServer) registerAPIGroup(subPrefix string, groupRegistrar func(router fiber.Router)) {
	if !strings.HasPrefix(subPrefix, "/") && subPrefix != "" {
		subPrefix = "/" + subPrefix
	}
	fullPrefix := s.basePath + subPrefix
	if s.basePath == "/" && subPrefix == "" {
		fullPrefix = "/"
	} else if s.basePath == "/" && strings.HasPrefix(subPrefix, "/") {
		fullPrefix = subPrefix
	} else if s.basePath != "" && subPrefix == "" {
		fullPrefix = s.basePath
	}

	group := s.apiBaseRouter.Group(subPrefix)
	groupRegistrar(group)
	s.log.Info("registered API group", slog.String("fullPrefix", fullPrefix))
}

func validateArguments(cfg core.AppConfigProvider, logger *slog.Logger, baseApiPrefix *string) {
	if cfg == nil {
		log.Fatal("server configuration is missing for HTTP server initialization")
	}
	if logger == nil {
		log.Fatal("logger instance is missing for HTTP server initialization")
	}
	if baseApiPrefix != nil {
		if *baseApiPrefix == "" {
			logger.Warn("baseApiPrefix is empty, API routes will be registered at the root.")
		} else if !strings.HasPrefix(*baseApiPrefix, "/") {
			*baseApiPrefix = "/" + *baseApiPrefix
			logger.Warn("baseApiPrefix did not start with '/', prepended it.", slog.String("newPrefix", *baseApiPrefix))
		}
	}
}
