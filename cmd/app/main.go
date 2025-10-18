package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zhunismp/intent-products-api/internal/interfaces/db"
	"github.com/zhunismp/intent-products-api/internal/interfaces/http/handlers"
	"github.com/zhunismp/intent-products-api/internal/applications/services"
	"github.com/zhunismp/intent-products-api/internal/infrastructure/client"
	"github.com/zhunismp/intent-products-api/internal/infrastructure/config"
)

func main() {
	cfg := config.Load()
	log.Printf("🚀 Starting %s in %s mode...", cfg.AppName, cfg.Env) // TODO: replace with cool banner.

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)

	// Initialize dependencies
	mongoClient := client.NewMongoClient(ctx, cfg)
	productRepository := db.NewProductRepositoryImpl(ctx, mongoClient.Database(cfg.Mongo.Database))
	productService := services.NewProductService(productRepository)
	productHttpHandler := handlers.NewProductHttpHandler(productService)

	// Initialize router
	r := chi.NewRouter()
	initalizeRoutes(r, productHttpHandler)

	// Start HTTP server
	srv := startHttpServer(ctx, cfg, r)
	gracefulShutdown(ctx, srv, cancel)

	defer func() {
		log.Print("Shutting down...")
		mongoClient.Disconnect(ctx)
		cancel()
	}()
}

func startHttpServer(ctx context.Context, cfg *config.Config, router *chi.Mux) *http.Server {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTPPort),
		Handler: router,
	}

	go func() {
		log.Printf("🚀 HTTP server running on port %d", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ ListenAndServe: %v", err)
		}
	}()

	return srv
}

func gracefulShutdown(ctx context.Context, srv *http.Server, cancelFunc context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("⚡ Received signal: %v. Shutting down...", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Print("✅ Server exited gracefully")
	}

	cancelFunc()
}

func initalizeRoutes(r chi.Router, productHttpHandler *handlers.ProductHttpHandler) {

	// product routes
	r.Route("/api/v1/products", func(r chi.Router) {
		r.Post("/", productHttpHandler.CreateProduct)
	})

	// health check routes
	r.Route("/api/v1/health", func(r chi.Router) {
		r.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})
	})
}
