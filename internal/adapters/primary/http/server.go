package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	cors "github.com/gofiber/fiber/v3/middleware/cors"
	limiter "github.com/gofiber/fiber/v3/middleware/limiter"
	logger "github.com/gofiber/fiber/v3/middleware/logger"
	recover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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

	app.Use(requestid.New())

	// Use slog via Done callback
	app.Use(logger.New(logger.Config{
		Stream: io.Discard,
		Done: func(c fiber.Ctx, logString []byte) {
			slogLogger.Info("http request",
				slog.String("request_id", requestid.FromContext(c)),
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.Int("status", c.Response().StatusCode()),
				slog.String("ip", c.IP()),
			)
		},
	}))

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
				"errorMessage": "Too many requests, please try again later.",
			})
		},
	}))

	apiGroup := app.Group(baseApiPrefix)
	slogLogger.Info("Fiber HTTP server core initialized with middleware.", slog.String("baseApiPrefix", baseApiPrefix))

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
	s.log.Info("Attempting to start HTTP server...", slog.String("address", serverAddr))

	go func() {
		if err := s.fiberApp.Listen(serverAddr); err != nil && err.Error() != "http: Server closed" {
			s.log.Error("Failed to start HTTP server listener")
		}
	}()

	s.log.Info("HTTP server listener started.", slog.String("address", serverAddr))
}

func (s *HttpServer) SetupRoute(routeGroup *RouteGroup) {
	if routeGroup.product == nil {
		s.log.Error("Failed to set up route")
	}

	productHandler := routeGroup.product

	s.registerAPIGroup("/health", func(router fiber.Router) {
		router.Get("/liveness", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"message": "application is running"}) })
	})

	s.registerAPIGroup("/products", func(router fiber.Router) {
		// core product
		router.Get("/:id", productHandler.GetProduct)
		router.Get("/status/:status", productHandler.GetProductByStatus)
		router.Post("/", productHandler.CreateProduct)
		router.Put("/position", productHandler.MoveProductPosition)
		router.Delete("/:id", productHandler.DeleteProduct)
	})
}

func (s *HttpServer) GracefulShutdown() {
	s.log.Info("Gracefully shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(shutdownCtx); err != nil {
		s.log.Error("Error during server shutdown")
	} else {
		s.log.Info("HTTP server shutdown gracefully.")
	}

	s.log.Info("Cleanup finished. Exiting.")
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
	s.log.Info("Registered API group", slog.String("fullPrefix", fullPrefix))
}

func validateArguments(cfg core.AppConfigProvider, otelLogger *slog.Logger, baseApiPrefix *string) {
	if cfg == nil {
		log.Fatal("Server configuration is missing for HTTP server initialization")
	}
	if otelLogger == nil {
		log.Fatal("Logger instance is missing for HTTP server initialization")
	}
	if baseApiPrefix != nil {
		if *baseApiPrefix == "" {
			otelLogger.Warn("baseApiPrefix is empty, API routes will be registered at the root.")
		} else if !strings.HasPrefix(*baseApiPrefix, "/") {
			*baseApiPrefix = "/" + *baseApiPrefix
			otelLogger.Warn("baseApiPrefix did not start with '/', prepended it.", slog.String("newPrefix", *baseApiPrefix))
		}
	}
}
